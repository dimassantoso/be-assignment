package resthandler

import (
	"billing-engine/internal/modules/auth/domain"
	"billing-engine/pkg/shared/usecase"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/codebase/factory/dependency"
	"github.com/golangid/candi/codebase/interfaces"
	"github.com/golangid/candi/tracer"
	"github.com/golangid/candi/wrapper"
)

// RestHandler handler
type RestHandler struct {
	mw        interfaces.Middleware
	uc        usecase.Usecase
	validator interfaces.Validator
}

// NewRestHandler create new rest handler
func NewRestHandler(uc usecase.Usecase, deps dependency.Dependency) *RestHandler {
	return &RestHandler{
		uc: uc, mw: deps.GetMiddleware(), validator: deps.GetValidator(),
	}
}

// Mount handler with root "/"
// handling version in here
func (h *RestHandler) Mount(root interfaces.RESTRouter) {
	v1Auth := root.Group(candihelper.V1 + "/auth")

	v1Auth.POST("/login", h.login, h.mw.HTTPBasicAuth)
	v1Auth.POST("", h.createAuth, h.mw.HTTPBearerAuth)
}

// CreateAuth documentation
// @Summary			Create Auth
// @Description		API for create auth
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			data	body	domain.RequestAuth	true	"Body Data"
// @Success			200	{object}	error
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/auth [post]
func (h *RestHandler) createAuth(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "AuthDeliveryREST:CreateAuth")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("auth/save", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestAuth
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	err := h.uc.Auth().CreateAuth(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusCreated, "Success").JSON(rw)
}

// CreateAuth documentation
// @Summary			Create Login
// @Description		API for login
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			data	body	domain.RequestLogin	true	"Body Data"
// @Success			200	{object}	domain.ResponseLogin
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		BasicAuth
// @Router			/v1/auth/login [post]
func (h *RestHandler) login(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "AuthDeliveryREST:Login")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("auth/login", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestLogin
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	result, attemptLogin, err := h.uc.Auth().LoginAuth(ctx, &payload)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "password not match" {
			if attemptLogin > 0 {
				wrapper.NewHTTPResponse(http.StatusForbidden, fmt.Sprintf("Email, username, or password is invalid. You have %d attempts left.", attemptLogin)).JSON(rw)
				return
			}

			wrapper.NewHTTPResponse(http.StatusUnprocessableEntity, "You have reached sign in attempt limit. Please try again after 5 minutes.").JSON(rw)
			return
		}
		wrapper.NewHTTPResponse(http.StatusBadRequest, "failed to login").JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success", result).JSON(rw)
}
