package task_service

import (
	"context"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

func (s *TaskService) GetTask(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error getting task from repository: %w", err)
	}
	return task, nil
}
