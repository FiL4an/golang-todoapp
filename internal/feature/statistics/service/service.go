package statistics_service

import (
	"context"
	"time"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	staticticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		staticticsRepository: statisticsRepository,
	}
}
