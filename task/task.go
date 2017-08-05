package task

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/node"
	"time"
)

var (
	registry = []*Task{}
)

type Task struct {
	Name    string
	Hours   []int
	Mins    []int
	Retry   bool
	Handler func(*database.Database) error
}

func (t *Task) scheduled(hour, min int) bool {
	for _, h := range t.Hours {
		if h == hour {
			for _, m := range t.Mins {
				if m == min {
					return true
				}
			}
		}
	}
	return false
}

func (t *Task) run(now time.Time) {
	db := database.GetDatabase()
	defer db.Close()

	job := &Job{
		Id: fmt.Sprintf(
			"%s-%d", t.Name, now.Unix()-int64(now.Second())),
		Name:      t.Name,
		State:     Running,
		Retry:     t.Retry,
		NodeId:    node.Self.Id,
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
		job.Failed(db)
		return
	}

	job.Finished(db)
}

func runScheduler() {
	now := time.Now()
	curHour := now.Hour()
	curMin := now.Minute()

	for {
		time.Sleep(1 * time.Second)

		now = time.Now()
		hour := now.Hour()
		min := now.Minute()

		if curHour == hour && curMin == min {
			continue
		}
		curHour = hour
		curMin = min

		for _, task := range registry {
			if task.scheduled(hour, min) {
				go task.run(now)
			}
		}
	}
}

func register(task *Task) {
	registry = append(registry, task)
}

func Init() {
	go runScheduler()
}
