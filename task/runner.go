package task

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/v2/bson"
	"github.com/pritunl/mongo-go-driver/v2/mongo/options"
	"github.com/pritunl/pritunl-zero/constants"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/version"
	"github.com/sirupsen/logrus"
)

func runnerCheck() time.Duration {
	return time.Duration(settings.System.TaskRunnerCheck) * time.Second
}

func runnerHeartbeat() time.Duration {
	return time.Duration(settings.System.TaskRunnerHeartbeat) * time.Second
}

func runnerTtl() time.Duration {
	return time.Duration(settings.System.TaskRunnerTtl) * time.Second
}

func runnerCooldown() time.Duration {
	return time.Duration(settings.System.TaskRunnerCooldown) * time.Second
}

func (t *Task) runnerSlotId(slot int) string {
	return fmt.Sprintf("%s-runner-%d", t.Name, slot)
}

func (t *Task) runnerWorkers() (workers int) {
	if t.Workers != nil {
		workers = t.Workers()
	}
	if workers < 1 {
		workers = 1
	}

	return
}

func (t *Task) runnerFindSlots(db *database.Database, now time.Time) (
	slots []int, err error) {

	coll := db.Tasks()
	slots = []int{}

	slotIds := []string{}
	for i := 0; i < t.runnerWorkers(); i++ {
		slotIds = append(slotIds, t.runnerSlotId(i))
	}

	cursor, err := coll.Find(db, &bson.M{
		"_id": &bson.M{
			"$in": slotIds,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}
	defer cursor.Close(db)

	active := map[string]bool{}
	for cursor.Next(db) {
		job := &Job{}
		err = cursor.Decode(job)
		if err != nil {
			err = database.ParseError(err)
			return
		}

		if job.State != Finished && now.Sub(job.Timestamp) < runnerTtl() {
			active[job.Id] = true
		}
	}

	err = cursor.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	for i, slotId := range slotIds {
		if !active[slotId] {
			slots = append(slots, i)
		}
	}

	return
}

func (t *Task) runnerReserve(db *database.Database, slot int) (
	job *Job, err error) {

	coll := db.Tasks()
	now := time.Now()

	newJob := &Job{
		Id:        t.runnerSlotId(slot),
		Name:      t.Name,
		Type:      Runner,
		State:     Running,
		Node:      node.Self.Id,
		Token:     bson.NewObjectID(),
		Started:   now,
		Timestamp: now,
	}

	_, err = coll.UpdateOne(db, &bson.M{
		"_id": newJob.Id,
		"$or": []*bson.M{
			&bson.M{
				"state": Finished,
			},
			&bson.M{
				"timestamp": &bson.M{
					"$lt": now.Add(-runnerTtl()),
				},
			},
		},
	}, &bson.M{
		"$set": &bson.M{
			"name":      newJob.Name,
			"type":      newJob.Type,
			"state":     newJob.State,
			"node":      newJob.Node,
			"token":     newJob.Token,
			"started":   newJob.Started,
			"timestamp": newJob.Timestamp,
		},
	}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		err = database.ParseError(err)
		if _, ok := err.(*database.DuplicateKeyError); ok {
			err = nil
		}
		return
	}

	job = newJob
	return
}

func (t *Task) runnerHold(holdCtx context.Context, job *Job,
	cancel context.CancelCauseFunc) {

	db := database.GetDatabase()
	defer db.Close()

	coll := db.Tasks()
	lastUpdate := time.Now()
	retryAfter := time.Time{}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-holdCtx.Done():
			return
		case <-ticker.C:
		}

		if constants.Interrupt {
			t.runnerRelease(db, job, Finished)
			cancel(&ShutdownError{
				errors.New("task: Runner shutdown"),
			})
			return
		}

		if time.Since(lastUpdate) < runnerHeartbeat() ||
			time.Now().Before(retryAfter) {

			continue
		}

		resp, err := coll.UpdateOne(db, &bson.M{
			"_id":   job.Id,
			"token": job.Token,
		}, &bson.M{
			"$set": &bson.M{
				"timestamp": time.Now(),
			},
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"task":  t.Name,
				"slot":  job.Id,
				"error": database.ParseError(err),
			}).Error("task: Failed to update runner reservation")

			if time.Since(lastUpdate) >= runnerTtl() {
				cancel(&ReservationLostError{
					errors.New("task: Runner reservation lost"),
				})
				return
			}

			retryAfter = time.Now().Add(5 * time.Second)
			continue
		}

		if resp.MatchedCount == 0 {
			cancel(&ReservationLostError{
				errors.New("task: Runner reservation lost"),
			})
			return
		}

		lastUpdate = time.Now()
	}
}

func (t *Task) runnerRelease(db *database.Database, job *Job,
	state string) {

	coll := db.Tasks()

	_, err := coll.UpdateOne(db, &bson.M{
		"_id":   job.Id,
		"token": job.Token,
	}, &bson.M{
		"$set": &bson.M{
			"state":     state,
			"timestamp": time.Now(),
		},
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"task":  t.Name,
			"slot":  job.Id,
			"error": database.ParseError(err),
		}).Error("task: Failed to release runner reservation")
	}
}

func (t *Task) runnerHandle(ctx context.Context,
	db *database.Database) (err error) {

	defer func() {
		panc := recover()
		if panc != nil {
			logrus.WithFields(logrus.Fields{
				"task":  t.Name,
				"trace": string(debug.Stack()),
				"panic": panc,
			}).Error("task: Panic in runner task")

			err = fmt.Errorf("task: Panic in runner task")
		}
	}()

	err = t.RunnerHandler(ctx, db)
	return
}

func (t *Task) runnerRun(db *database.Database, job *Job) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	timer := time.AfterFunc(t.Duration, func() {
		cancel(&DurationElapsedError{
			errors.New("task: Runner duration elapsed"),
		})
	})
	defer timer.Stop()

	holdCtx, holdCancel := context.WithCancel(context.Background())
	holdDone := make(chan struct{})
	go func() {
		defer close(holdDone)
		t.runnerHold(holdCtx, job, cancel)
	}()

	err := t.runnerHandle(ctx, db)

	holdCancel()
	<-holdDone

	cause := context.Cause(ctx)
	_, lost := cause.(*ReservationLostError)

	if _, ok := cause.(*ShutdownError); ok {
		return
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"task":  t.Name,
			"slot":  job.Id,
			"error": err,
		}).Error("task: Runner task failed")
	}

	if lost {
		logrus.WithFields(logrus.Fields{
			"task": t.Name,
			"slot": job.Id,
		}).Warn("task: Runner task lost reservation")
		return
	}

	if err != nil {
		t.runnerRelease(db, job, Failed)
	} else {
		t.runnerRelease(db, job, Finished)
	}
}

func (t *Task) runRunner() {
	go func() {
		defer func() {
			panc := recover()
			if panc != nil {
				logrus.WithFields(logrus.Fields{
					"task":  t.Name,
					"trace": string(debug.Stack()),
					"panic": panc,
				}).Error("task: Panic in runner check")
			}
		}()

		if t.DebugNodes != nil {
			matched := false
			for _, ndeName := range t.DebugNodes {
				if node.Self.Name == ndeName {
					matched = true
				}
			}
			if !matched {
				return
			}
		}

		if constants.Interrupt {
			return
		}

		if !t.running.CompareAndSwap(0, time.Now().UnixNano()) {
			return
		}
		defer t.running.Store(0)

		if time.Now().UnixNano() < t.cooldown.Load() {
			return
		}

		time.Sleep(time.Duration(utils.RandInt(0, 1000)) * time.Millisecond)

		db := database.GetDatabase()
		defer db.Close()

		slots, err := t.runnerFindSlots(db, time.Now())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"task":  t.Name,
				"error": err,
			}).Error("task: Runner check failed")
			return
		}
		if len(slots) == 0 {
			return
		}

		if t.Version != 0 {
			supported, err := version.Check(db, t.Name, t.Version)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"task":  t.Name,
					"error": err,
				}).Error("task: Version check failed")
				return
			}

			if !supported {
				return
			}
		}

		if constants.Interrupt {
			return
		}

		var job *Job
		for _, slot := range slots {
			job, err = t.runnerReserve(db, slot)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"task":  t.Name,
					"error": err,
				}).Error("task: Runner reserve failed")
				return
			}
			if job != nil {
				break
			}
		}
		if job == nil {
			return
		}

		t.runnerRun(db, job)

		t.cooldown.Store(time.Now().Add(runnerCooldown()).UnixNano())
	}()
}

func Sleep(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
