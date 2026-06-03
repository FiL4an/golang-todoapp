package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_errors "github.com/FiL4an/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userID *int, from *time.Time, to *time.Time) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to ` must be after `from`: %w ", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.staticticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w ", err)
	}
	statistics := calcStatistics(tasks)

	return statistics, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{
			TasksCreated:               0,
			TasksCompleted:             0,
			TasksCompletedRate:         nil,
			TasksAverageCompletionTime: nil,
		}
	}

	taskCreated := len(tasks)
	taskscompleted := 0

	var totatlComplitionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			taskscompleted++
		}
		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totatlComplitionDuration += *completionDuration
		}
	}

	tasksCompleteRate := float64(taskscompleted) / float64(taskCreated) * 100
	var tasksAVGCompletionTime *time.Duration

	if taskscompleted > 0 && totatlComplitionDuration != 0 {
		avg := totatlComplitionDuration / time.Duration(taskscompleted)
		tasksAVGCompletionTime = &avg
	}

	return domain.Statistics{
		TasksCreated:               taskCreated,
		TasksCompleted:             taskscompleted,
		TasksCompletedRate:         &tasksCompleteRate,
		TasksAverageCompletionTime: tasksAVGCompletionTime,
	}
}
