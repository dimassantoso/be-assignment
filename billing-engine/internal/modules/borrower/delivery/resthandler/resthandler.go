package resthandler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"billing-engine/internal/modules/borrower/domain"
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
	v1Borrower := root.Group(candihelper.V1+"/borrower", h.mw.HTTPBearerAuth)

	v1Borrower.GET("/", h.getAllBorrower)
	v1Borrower.GET("/:id", h.getDetailBorrowerByID)
	v1Borrower.POST("/", h.createBorrower)
	v1Borrower.PUT("/:id", h.updateBorrower)
	v1Borrower.DELETE("/:id", h.deleteBorrower)
	v1Borrower.GET("/:id/delinquent-check", h.delinquentCheck)
}

// GetAllBorrower documentation
// @Summary			Get All Borrower
// @Description		API for get all borrower
// @Tags			Borrower
// @Accept			json
// @Produce			json
// @Param			page	query	string	false	"Page with default value is 1"
// @Param			limit	query	string	false	"Limit with default value is 10"
// @Param			search	query	string	false	"Search"
// @Param			orderBy	query	string	false	"Order By"
// @Param			sort	query	string	false	"Sort (ASC DESC)"
// @Success			200	{object}	domain.ResponseBorrowerList
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/borrower [get]
func (h *RestHandler) getAllBorrower(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BorrowerDeliveryREST:GetAllBorrower")
	defer trace.Finish()

	var filter domain.FilterBorrower
	if err := candihelper.ParseFromQueryParam(req.URL.Query(), &filter); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed parse filter", err).JSON(rw)
		return
	}

	if err := h.validator.ValidateDocument("borrower/get_all", filter); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate filter", err).JSON(rw)
		return
	}

	result, err := h.uc.Borrower().GetAllBorrower(ctx, &filter)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	res := wrapper.NewHTTPResponse(http.StatusOK, "Success", result.Data)
	res.Meta = result.Meta
	res.JSON(rw)
}

// GetDetailBorrower documentation
// @Summary			Get Detail Borrower
// @Description		API for get detail borrower
// @Tags			Borrower
// @Accept			json
// @Produce			json
// @Param			id	path	string	true	"ID"
// @Success			200	{object}	domain.ResponseBorrower
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/borrower/{id} [get]
func (h *RestHandler) getDetailBorrowerByID(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BorrowerDeliveryREST:GetDetailBorrowerByID")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	data, err := h.uc.Borrower().GetDetailBorrower(ctx, id)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success", data).JSON(rw)
}

// CreateBorrower documentation
// @Summary			Create Borrower
// @Description		API for create borrower
// @Tags			Borrower
// @Accept			json
// @Produce			json
// @Param			data	body	domain.RequestBorrower	true	"Body Data"
// @Success			200	{object}	domain.ResponseBorrower
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/borrower [post]
func (h *RestHandler) createBorrower(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BorrowerDeliveryREST:CreateBorrower")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("borrower/save", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestBorrower
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	res, err := h.uc.Borrower().CreateBorrower(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusCreated, "Success", res).JSON(rw)
}

// UpdateBorrower documentation
// @Summary			Update Borrower
// @Description		API for update borrower
// @Tags			Borrower
// @Accept			json
// @Produce			json
// @Param			id	path	string	true	"ID"
// @Param			data	body	domain.RequestBorrower	true	"Body Data"
// @Success			200	{object}	domain.ResponseBorrower
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/borrower/{id} [put]
func (h *RestHandler) updateBorrower(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BorrowerDeliveryREST:UpdateBorrower")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	if err := h.validator.ValidateDocument("borrower/save", body); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed validate payload", err).JSON(rw)
		return
	}

	var payload domain.RequestBorrower
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	payload.ID, _ = strconv.Atoi(restserver.URLParam(req, "id"))
	err := h.uc.Borrower().UpdateBorrower(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success").JSON(rw)
}

// DeleteBorrower documentation
// @Summary			Delete Borrower
// @Description		API for delete borrower
// @Tags			Borrower
// @Accept			json
// @Produce			json
// @Param			id	path	string	true	"ID"
// @Success			200	{object}	domain.ResponseBorrower
// @Success			400	{object}	wrapper.HTTPResponse
// @Security		ApiKeyAuth
// @Router			/v1/borrower/{id} [delete]
func (h *RestHandler) deleteBorrower(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BorrowerDeliveryREST:DeleteBorrower")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	if err := h.uc.Borrower().DeleteBorrower(ctx, id); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success").JSON(rw)
}

func (h *RestHandler) delinquentCheck(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "BorrowerDeliveryREST:DelinquentCheck")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	data, err := h.uc.Borrower().DelinquentCheck(ctx, id)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success", data).JSON(rw)
}
