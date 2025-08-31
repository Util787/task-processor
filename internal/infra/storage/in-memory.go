package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Util787/task-processor/internal/common"
	"github.com/Util787/task-processor/internal/models"
)

const newMapCap = 100

type InMemoryTaskStateStorage struct {
	tasks map[string]models.TaskState
	rwmu  *sync.RWMutex
}

func NewInMemoryTaskStateStorage() *InMemoryTaskStateStorage {
	return &InMemoryTaskStateStorage{
		tasks: make(map[string]models.TaskState, newMapCap),
		rwmu:  &sync.RWMutex{},
	}
}

func (s *InMemoryTaskStateStorage) GetTaskState(ctx context.Context, taskID string) (models.TaskState, error) {
	op := common.GetOperationName()

	if ctx.Err() != nil {
		return "", fmt.Errorf("%s: %w", op, ctx.Err())
	}

	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	state, ok := s.tasks[taskID]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, models.ErrTaskNotFound)
	}
	return state, nil
}

func (s *InMemoryTaskStateStorage) SetTaskState(ctx context.Context, taskID string, state models.TaskState) error {
	op := common.GetOperationName()

	if ctx.Err() != nil {
		return fmt.Errorf("%s: %w", op, ctx.Err())
	}

	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	s.tasks[taskID] = state
	return nil
}
