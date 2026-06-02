package task_transport_http

import (
	"net/http"

	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_request "github.com/FiL4an/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/FiL4an/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

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
