package planning

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sort"
	"strings"
	"time"
)

type TaskState string

const (
	Queued  TaskState = "queued"
	Running TaskState = "running"
	Blocked TaskState = "blocked"
	Done    TaskState = "done"
	Failed  TaskState = "failed"
)

type Task struct {
	ID, Name, Zone        string
	State                 TaskState
	Priority              int
	Dependencies          []string
	StartedAt, FinishedAt *time.Time
	Attempts              int
}

func (t Task) Valid() error {
	if t.ID == "" || strings.TrimSpace(t.Name) == "" || t.Zone == "" || t.Priority < 1 || t.Priority > 5 {
		return common.ErrInvalid
	}
	return nil
}
func (t *Task) Start(now time.Time) error {
	if e := t.Valid(); e != nil {
		return e
	}
	if t.State != Queued {
		return fmt.Errorf("%w: task not queued", common.ErrConflict)
	}
	t.State = Running
	t.StartedAt = &now
	t.Attempts++
	return nil
}
func (t *Task) Finish(now time.Time) error {
	if t.State != Running {
		return common.ErrConflict
	}
	t.State = Done
	t.FinishedAt = &now
	return nil
}
func (t *Task) Fail(now time.Time) error {
	if t.State != Running {
		return common.ErrConflict
	}
	t.State = Failed
	t.FinishedAt = &now
	return nil
}
func (t Task) Ready(done map[string]bool) bool {
	for _, d := range t.Dependencies {
		if !done[d] {
			return false
		}
	}
	return true
}
func Schedule(tasks []Task) []Task {
	out := append([]Task(nil), tasks...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].ID < out[j].ID
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}
