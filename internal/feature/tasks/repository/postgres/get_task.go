package task_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_errors "github.com/FiL4an/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/FiL4an/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at,author_user_id
	FROM todoapp.task
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var task TaskModel

	err := row.Scan(
		&task.ID,
		&task.Version,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.AuthorUserID)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id= '%d': %w", id, core_errors.ErrNotFound)
		}
		return domain.Task{}, fmt.Errorf("Scan error : %w", err)
	}

	taskDomain := taskDomainFromModel(task)

	return taskDomain, nil
}
