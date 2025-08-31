package usecase

import (
	"context"

	"github.com/Util787/task-processor/internal/models"
)

type TaskStateStorage interface {
	GetTaskState(ctx context.Context, taskID string) (models.TaskState, error)
}

type TaskProcessQueue interface {
	EnqueueTask(ctx context.Context, task models.Task) error
}

type TaskUsecase struct {
	taskStorage      TaskStateStorage
	taskProcessQueue TaskProcessQueue
}

func NewTaskUsecase(taskStorage TaskStateStorage, taskProcessQueue TaskProcessQueue) *TaskUsecase {
	return &TaskUsecase{
		taskStorage:      taskStorage,
		taskProcessQueue: taskProcessQueue,
	}
}
