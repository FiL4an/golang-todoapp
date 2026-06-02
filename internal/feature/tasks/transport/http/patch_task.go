package task_transport_http

import (
	"fmt"
	"net/http"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_request "github.com/FiL4an/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/FiL4an/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/FiL4an/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`title` can't be null")
		}
		titltLen := len([]rune(*r.Title.Value))
		if titltLen < 1 || titltLen > 100 {
			return fmt.Errorf("`title` length must be between 1 and 100 characters")
		}
	}
	if r.Description.Set {
		if r.Description.Value != nil {
			descLen := len([]rune(*r.Description.Value))
			if descLen < 1 || descLen > 1000 {
				return fmt.Errorf("`description` length must be between 1 and 1000 characters")
			}
		}

	}
	if !r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`completed` can't be null")
		}

	}
	return nil
}

func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}
	var request PatchTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}
	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.taskService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}
	jsonResponse := PatchTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(jsonResponse, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
