package task_postgres_repository

import (
	"context"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title,description,completed, created_at , completed_at,author_user_id 
	FROM todoapp.task
	%s
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2
	`

	args := []any{limit, offset}
	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskmodels []TaskModel
	for rows.Next() {
		var taskmodel TaskModel
		if err := rows.Scan(
			&taskmodel.ID,
			&taskmodel.Version,
			&taskmodel.Title,
			&taskmodel.Description,
			&taskmodel.Completed,
			&taskmodel.CreatedAt,
			&taskmodel.CompletedAt,
			&taskmodel.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("Scan tasks:%w", err)
		}
		taskmodels = append(taskmodels, taskmodel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan tasks: %w", err)
	}
	tasksDomains := taskDomainsFromModels(taskmodels)

	return tasksDomains, nil
}
