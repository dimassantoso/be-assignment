package resthandler

import (
	"billing-engine/internal/modules/billing/domain"
	"billing-engine/pkg/shared/usecase"
	"encoding/json"
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
	v1Billing := root.Group(candihelper.V1+"/billing", h.mw.HTTPBearerAuth)

	v1Billing.POST("/repayment", h.createBilling)
}

// CreateBilling documentation
// @Summary			Create Billing
// @Description		API for create billing
// @Tags			Billing
// @Accept			json
// @Produce			json
// @Param			data	body	domain.RequestBilling	true	"Body Data"
// @Success			200	{object}	domain.ResponseBilling
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/billing [post]
func (h *RestHandler) createBilling(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BillingDeliveryREST:CreateBilling")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("billing/save", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestBillingRepayment
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	err := h.uc.Billing().BillingRepayment(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success").JSON(rw)
}
