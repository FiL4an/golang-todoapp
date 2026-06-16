package task_transport_http

import (
	"net/http"

	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_request "github.com/FiL4an/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/FiL4an/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary Получть задачу
// @Description Получение задачи по author ID
// @Tags tasks
// @Produce json
// @Param id path int true  "ID получаемой задачи"
// @Success 200 {object}  GetTaskResponse "Получение задачи автора"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Not found task"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHanlder := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHanlder.ErrorResponse(err,
			"failed to get taskID path value")
		return
	}
	taskDomain, err := h.taskService.GetTask(ctx, userID)
	if err != nil {
		responseHanlder.ErrorResponse(err, "failed to get task")
		return
	}
	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHanlder.JSONResponse(response, http.StatusOK)
}
