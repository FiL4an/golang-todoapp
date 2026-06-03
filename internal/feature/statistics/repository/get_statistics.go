package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
		SELECT id , version, title, description, completed, created_at,completed_at, author_user_id 
		FROM todoapp.task
	`)

	args := []any{}
	condition := []string{}
	if userID != nil {
		condition = append(condition, fmt.Sprintf("author_user_id=%d", len(args)+1))
		args = append(args, userID)
	}
	if from != nil {
		condition = append(condition, fmt.Sprintf("created_at>=%d", len(args)+1))
		args = append(args, from)
	}
	if to != nil {
		condition = append(condition, fmt.Sprintf("created_at<%d", len(args)+1))
		args = append(args, to)
	}
	if len(condition) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(condition, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select task:%w ", err)
	}
	defer rows.Close()

	var tasksModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		if err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("Scan error:%w", err)
		}
		tasksModels = append(tasksModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows :%w ", err)
	}

	domainsTasks := taskDomainsFromModels(tasksModels)
	return domainsTasks, nil
}
