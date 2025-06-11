package task

import (
	"fmt"
	"runtime/debug"
	"slices"
	"time"

	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/utils"
	"github.com/pritunl/pritunl-zero/version"
	"github.com/sirupsen/logrus"
)

var (
	registry = []*Task{}
)

type Task struct {
	Name       string
	Version    int
	Hours      []int
	Minutes    []int
	Seconds    time.Duration
	Retry      bool
	Handler    func(*database.Database) error
	RunOnStart bool
	Local      bool
	DebugNodes []string
	timestamp  time.Time
}

func (t *Task) scheduled(hour, min int) bool {
	for _, h := range t.Hours {
		if h == hour {
			if slices.Contains(t.Minutes, min) {
				return true
			}
		}
	}
	return false
}

func (t *Task) runShared(db *database.Database, now time.Time) {
	defer func() {
		panc := recover()
		if panc != nil {
			logrus.WithFields(logrus.Fields{
				"trace": string(debug.Stack()),
				"panic": panc,
			}).Error("sync: Panic in run task")
		}
	}()

	if t.Seconds == 0 {
		time.Sleep(time.Duration(utils.RandInt(0, 1000)) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(utils.RandInt(0, 300)) * time.Millisecond)
	}

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

	id := fmt.Sprintf("%s-%d", t.Name, now.Unix()-int64(now.Second()))
	if t.Seconds != 0 {
		id += fmt.Sprintf("-%d", GetBlock(now, t.Seconds))
	}

	job := &Job{
		Id:        id,
		Name:      t.Name,
		State:     Running,
		Retry:     t.Retry,
		Node:      node.Self.Id,
		Timestamp: time.Now(),
	}

	reserved, err := job.Reserve(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"task":  t.Name,
			"error": err,
		}).Error("task: Task reserve failed")
		return
	}

	if !reserved {
		return
	}

	err = t.Handler(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"task":  t.Name,
			"error": err,
		}).Error("task: Task failed")
		_ = job.Failed(db)
		return
	}

	_ = job.Finished(db)
}

func (t *Task) runLocal(db *database.Database, now time.Time) {
	defer func() {
		panc := recover()
		if panc != nil {
			logrus.WithFields(logrus.Fields{
				"trace": string(debug.Stack()),
				"panic": panc,
			}).Error("sync: Panic in run local task")
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

	id := fmt.Sprintf("%s-%d", t.Name, now.Unix()-int64(now.Second()))
	if t.Seconds != 0 {
		id += fmt.Sprintf("-%d", GetBlock(now, t.Seconds))
	}

	err := t.Handler(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"task":  t.Name,
			"error": err,
		}).Error("task: Local task failed")
		return
	}
}

func (t *Task) run(now time.Time) {
	go func() {
		db := database.GetDatabase()
		defer db.Close()

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
				logrus.WithFields(logrus.Fields{
					"task":    t.Name,
					"version": t.Version,
				}).Info("task: Skipping incompatible task")
				return
			}
		}

		curTimestamp := t.timestamp
		if !curTimestamp.IsZero() {
			if time.Since(curTimestamp) > 10*time.Minute {
				logrus.WithFields(logrus.Fields{
					"task_name": t.Name,
					"runtime":   time.Since(curTimestamp),
				}).Error("task: Task stuck running")
			}
			return
		}
		t.timestamp = time.Now()
		defer func() {
			t.timestamp = time.Time{}
		}()

		if t.Local {
			t.runLocal(db, now)
		} else {
			t.runShared(db, now)
		}
	}()
}

func runScheduler() {
	now := time.Now()
	curHour := now.Hour()
	curMin := now.Minute()
	curSecBlocks := map[time.Duration]int{}

	for _, task := range registry {
		if task.Seconds != 0 {
			curSecBlocks[task.Seconds] = GetBlock(now, task.Seconds)
		}

		if task.RunOnStart {
			go task.run(now)
		}
	}

	for {
		time.Sleep(1 * time.Second)

		now = time.Now()
		hour := now.Hour()
		min := now.Minute()

		for block, curSecBlock := range curSecBlocks {
			secBlock := GetBlock(now, block)

			if curSecBlock != secBlock {
				for _, task := range registry {
					if task.Seconds != 0 && task.Seconds == block &&
						task.scheduled(hour, min) {

						task.run(now)
					}
				}
			}

			curSecBlocks[block] = secBlock
		}

		if curHour == hour && curMin == min {
			continue
		}
		curHour = hour
		curMin = min

		for _, task := range registry {
			if task.Seconds == 0 && task.scheduled(hour, min) {
				task.run(now)
			}
		}
	}
}

func register(task *Task) {
	registry = append(registry, task)
}

func Init() (err error) {
	for _, task := range registry {
		if task.Version == 0 {
			continue
		}

		err = version.Set(database.GetDatabase(), task.Name, task.Version)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"task":    task.Name,
				"version": task.Version,
				"error":   err,
			}).Error("task: Failed to set task version")
			return
		}
	}

	go runScheduler()

	return
}

func GetBlock(n time.Time, d time.Duration) int {
	s := int(d.Seconds())
	return (n.Second() / s) * s
}
