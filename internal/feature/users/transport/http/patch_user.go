package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_request "github.com/FiL4an/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/FiL4an/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/FiL4an/golang-todoapp/internal/core/transport/http/types"
)

type PacthUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Максим Максимович"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887777"`
}

func (r *PacthUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL  ")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName must be between 3 and 100  symbols`")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))

			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber must be between 10 and 15  symbols`")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must be startswith '+' symbol")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary Изменение пользователя
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"phone_number":"+71112223344"`- устанавливает новый номер телефона в БД
// @Description 3. **Передан null**: `"phone_number":"null"`- очищает поле в БД (set to NULL)
// @Description  Ограничения: `full_name` не может быть выставлен как null
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемого пользователя"
// @Param request body PacthUserRequest true "PatchUser тело запроса"
// @Success 200 {object} PatchUserResponse "Успешно измененный пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	var request PacthUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err,
			"failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(fmt.Sprintf("PacthUserRequest fields:\nFullName: '%v'\nPhoneNumber: '%v'", request.FullName, request.PhoneNumber))
	w.WriteHeader(http.StatusOK)
}

func userPatchFromRequest(request PacthUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)

}
