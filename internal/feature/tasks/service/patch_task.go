package task_service

import (
	"context"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

func (s *TaskService) PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply patch to task: %w", err)
	}
	patchedtask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task in repository: %w", err)
	}
	return patchedtask, nil
}
