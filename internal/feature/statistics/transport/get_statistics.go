package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_request "github.com/FiL4an/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/FiL4an/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

// GetStatistics godoc
// @Summary Получение статистики
// @Description Получение статистики по задачам с опциаональной фильтрацией по user_id и/или врвеменному промежутку
// @Tags statistics
// @Produce json
// @Param user_id query int false "Филтрация статистики по конкретному пользователю"
// @Param from query  string false "Начало промежутка рассмотрения статстики(включительно)"
// @Param to query string false "Конец промежутка рассмотрения статистики(не включительно)"
// @Success 200 {object} GetStatisticsResponse "Успешное получение статистики"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (s *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)
	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/from/to query params")
		return
	}
	statistics, err := s.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)

}
func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParam      = "from"
		toQueryParam        = "to"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userID' query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}
	return userID, from, to, nil
}
