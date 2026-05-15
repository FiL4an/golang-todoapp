package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_request "github.com/FiL4an/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/FiL4an/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := h.getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit'/'offset' query param")
		return
	}

	usersDomain, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err,
			"failed to get users")
		return
	}

	response := GetUsersResponse(userDTOFromDomains(usersDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func (h *UsersHTTPHandler) getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return limit, offset, nil
}
