package taskprocessqueue

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Util787/task-processor/internal/common"
	"github.com/Util787/task-processor/internal/config"
	"github.com/Util787/task-processor/internal/models"
)

const defaultDelayMs = 100

type TaskStateStorage interface {
	SetTaskState(ctx context.Context, taskID string, state models.TaskState) error
}

type TaskProcessQueue struct {
	stateStorage TaskStateStorage
	taskChan     chan models.Task
	wg           *sync.WaitGroup
	log          *slog.Logger // the point of logger here is because task state only changing in this layer logger should be used only to log task state setting errors
	closed       atomic.Bool  // default is false
}

func NewTaskProcessQueue(ctx context.Context, log *slog.Logger, cfg config.TaskProcessQueueConfig, stateStorage TaskStateStorage) *TaskProcessQueue {
	q := &TaskProcessQueue{
		stateStorage: stateStorage,
		taskChan:     make(chan models.Task, cfg.QueueSize),
		wg:           &sync.WaitGroup{},
		log:          log,
	}

	for range cfg.Workers {
		go q.worker()
	}

	return q
}

// context shouldnt affect workers so just use context.Background()
func (q *TaskProcessQueue) worker() {
	log := q.log.With(slog.String("operation", common.GetOperationName()))

	for task := range q.taskChan {
		log = log.With(slog.String("task_id", task.Id))

		err := q.stateStorage.SetTaskState(context.Background(), task.Id, models.StateRunning)
		if err != nil {
			log.Warn("failed to set task state", slog.String("error", err.Error()))
		}

		if done := processTaskWithRetryImitation(&task); done {
			err := q.stateStorage.SetTaskState(context.Background(), task.Id, models.StateDone)
			if err != nil {
				log.Warn("failed to set task state", slog.String("error", err.Error()))
			}
		} else {
			err := q.stateStorage.SetTaskState(context.Background(), task.Id, models.StateFailed)
			if err != nil {
				log.Warn("failed to set task state", slog.String("error", err.Error()))
			}
		}
		q.wg.Done()
	}
}

func processTaskWithRetryImitation(task *models.Task) bool {
	for attempt := 1; attempt <= task.MaxRetries; attempt++ {
		workDuration := time.Duration(100+rand.Intn(401)) * time.Millisecond
		time.Sleep(workDuration)

		if rand.Float64() < 0.2 {
			if attempt < task.MaxRetries {
				backoff := time.Duration(defaultDelayMs<<attempt) * time.Millisecond
				jitter := time.Duration(rand.Intn(defaultDelayMs)) * time.Millisecond
				time.Sleep(backoff + jitter)
			}
			continue
		}
		return true
	}
	return false
}

func (q *TaskProcessQueue) EnqueueTask(ctx context.Context, task models.Task) error {
	op := common.GetOperationName()
	log := q.log.With(slog.String("operation", op), slog.String("task_id", task.Id))

	if ctx.Err() != nil {
		return fmt.Errorf("%s: %w", op, ctx.Err())
	}
	if q.closed.Load() {
		return fmt.Errorf("%s: queue is closed due to shutdown", op)
	}

	err := q.stateStorage.SetTaskState(ctx, task.Id, models.StateQueued)
	if err != nil {
		log.Warn("failed to set task state", slog.String("error", err.Error()))
	}
	q.taskChan <- task
	q.wg.Add(1)

	return nil
}

func (q *TaskProcessQueue) Shutdown() {
	q.closed.Store(true)
	close(q.taskChan)

	q.wg.Wait()
}
