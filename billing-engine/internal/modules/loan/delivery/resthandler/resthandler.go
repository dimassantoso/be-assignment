package resthandler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"billing-engine/internal/modules/loan/domain"
	"billing-engine/pkg/shared/usecase"

	"github.com/golangid/candi/candihelper"
	restserver "github.com/golangid/candi/codebase/app/rest_server"
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
	v1Loan := root.Group(candihelper.V1+"/loan", h.mw.HTTPBearerAuth)

	v1Loan.POST("/", h.createLoan)
	v1Loan.POST("/simulation", h.getLoanSimulation)
	v1Loan.GET("/outstanding/borrower/:borrower_id", h.getLoanOutstanding)
}

// CreateLoan documentation
// @Summary			Create Loan
// @Description		API for create loan
// @Tags			Loan
// @Accept			json
// @Produce			json
// @Param			data	body	domain.RequestLoan	true	"Body Data"
// @Success			200	{object}	domain.ResponseLoan
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		BearerToken
// @Router			/v1/loan [post]
func (h *RestHandler) createLoan(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "LoanDeliveryREST:CreateLoan")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("loan/save", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestLoan
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	res, err := h.uc.Loan().CreateLoan(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusCreated, "Success", res).JSON(rw)
}

// GetLoanSimulation documentation
// @Summary			Get Loan Simulation
// @Description		API for get loan simulation
// @Tags			Loan
// @Accept			json
// @Produce			json
// @Param			data	body	domain.RequestLoanSimulation	true	"Body Data"
// @Success			200	{object}	domain.ResponseLoanSimulation
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		BearerToken
// @Router			/v1/loan/simulation [post]
func (h *RestHandler) getLoanSimulation(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "LoanDeliveryREST:GetLoanSimulation")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("loan/simulation", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestLoanSimulation
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	res, err := h.uc.Loan().GetLoanSimulation(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success", res).JSON(rw)
}

// GetLoanSimulation documentation
// @Summary			Get Loan Simulation
// @Description		API for get loan simulation
// @Tags			Loan
// @Accept			json
// @Produce			json
// @Param			borrower_id	path	string	true	"borrowerID"
// @Success			200	{object}	domain.ResponseLoanOutstanding
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		BearerToken
// @Router			/v1/loan/simulation [post]
func (h *RestHandler) getLoanOutstanding(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "LoanDeliveryREST:GetOutstanding")
	defer trace.Finish()

	borrowerID, _ := strconv.Atoi(restserver.URLParam(req, "borrower_id"))
	res, err := h.uc.Loan().GetLoanOutstanding(ctx, borrowerID)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success", res).JSON(rw)
}
