package usecase

import (
	"context"
	"fmt"

	"github.com/Util787/task-processor/internal/common"
	"github.com/Util787/task-processor/internal/models"
)

func (u *TaskUsecase) EnqueueTask(ctx context.Context, task models.Task) error {
	op := common.GetOperationName()

	//validation
	err := validateTask(task)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, models.ErrValidation, err)
	}

	err = u.taskProcessQueue.EnqueueTask(ctx, task)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *TaskUsecase) GetTaskState(ctx context.Context, taskID string) (models.TaskState, error) {
	op := common.GetOperationName()

	state, err := u.taskStorage.GetTaskState(ctx, taskID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return state, nil
}

func validateTask(task models.Task) error {
	if task.Id == "" {
		return fmt.Errorf("task ID is required")
	}

	if task.MaxRetries <= 0 {
		return fmt.Errorf("max retries must be > 0")
	}

	return nil
}
