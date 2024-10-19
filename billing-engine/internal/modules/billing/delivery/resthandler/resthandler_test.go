package resthandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockusecase "billing-engine/pkg/mocks/modules/billing/usecase"
	mocksharedusecase "billing-engine/pkg/mocks/shared/usecase"

	"github.com/golangid/candi/candihelper"
	mockdeps "github.com/golangid/candi/mocks/codebase/factory/dependency"
	mockinterfaces "github.com/golangid/candi/mocks/codebase/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testCase struct {
	name, reqBody                       string
	wantValidateError, wantUsecaseError error
	wantRespCode                        int
}

var (
	errFoo = errors.New("Something error")
)

func TestNewRestHandler(t *testing.T) {
	mockMiddleware := &mockinterfaces.Middleware{}
	mockValidator := &mockinterfaces.Validator{}

	mockDeps := &mockdeps.Dependency{}
	mockDeps.On("GetMiddleware").Return(mockMiddleware)
	mockDeps.On("GetValidator").Return(mockValidator)

	handler := NewRestHandler(nil, mockDeps)
	assert.NotNil(t, handler)

	mockRoute := &mockinterfaces.RESTRouter{}
	mockRoute.On("Group", mock.Anything, mock.Anything).Return(mockRoute)
	mockRoute.On("POST", mock.Anything, mock.Anything, mock.Anything)
	handler.Mount(mockRoute)
}

func TestRestHandler_billingRepayment(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", reqBody: `{"billing_id":1,"payment_method_id":1}`, wantUsecaseError: nil, wantRespCode: http.StatusOK,
		},
		{
			name: "Testcase #2: Negative", reqBody: `{billing_id:1,"payment_method_id":1}`, wantUsecaseError: nil, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", reqBody: `{"billing_id":1,"payment_method_id":1}`, wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", reqBody: `{"billing_id":"a","payment_method_id":1}`, wantUsecaseError: nil, wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			billingUsecase := &mockusecase.BillingUsecase{}
			billingUsecase.On("BillingRepayment", mock.Anything, mock.Anything).Return(tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Billing").Return(billingUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.createBilling(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}
