package resthandler

import (
	"errors"
	"github.com/golangid/candi/candishared"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"billing-engine/internal/modules/loan/domain"
	mockusecase "billing-engine/pkg/mocks/modules/loan/usecase"
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
	mockRoute.On("GET", mock.Anything, mock.Anything, mock.Anything)
	mockRoute.On("POST", mock.Anything, mock.Anything, mock.Anything)
	handler.Mount(mockRoute)
}

func TestRestHandler_createLoan(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", reqBody: `{"borrower_id":1,"principle_amount":5000000,"duration_id":1}`, wantUsecaseError: nil, wantRespCode: http.StatusCreated,
		},
		{
			name: "Testcase #2: Negative", reqBody: `{"borrower_id":1,principle_amount:5000000,"duration_id":1}`, wantUsecaseError: nil, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", reqBody: `{"borrower_id":1,"principle_amount":5000000,"duration_id":1}`, wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", reqBody: `{"borrower_id":1,"principle_amount":5000000,"duration_id":1}`, wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			loanUsecase := &mockusecase.LoanUsecase{}
			loanUsecase.On("CreateLoan", mock.Anything, mock.Anything).Return(domain.ResponseLoan{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Loan").Return(loanUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.createLoan(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_getSimulation(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", reqBody: `{"principle_amount":5000000,"duration_id":1}`, wantUsecaseError: nil, wantRespCode: http.StatusOK,
		},
		{
			name: "Testcase #2: Negative", reqBody: `{principle_amount:5000000,"duration_id":1}`, wantUsecaseError: nil, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", reqBody: `{"principle_amount":5000000,"duration_id":1}`, wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", reqBody: `{"principle_amount":5000000,"duration_id":1}`, wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			loanUsecase := &mockusecase.LoanUsecase{}
			loanUsecase.On("GetLoanSimulation", mock.Anything, mock.Anything).Return([]domain.ResponseLoanSimulation{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Loan").Return(loanUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.getLoanSimulation(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_getLoanOutstanding(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", wantUsecaseError: nil, wantRespCode: http.StatusOK,
		},
		{
			name: "Testcase #2: Negative", wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			loanUsecase := &mockusecase.LoanUsecase{}
			loanUsecase.On("GetLoanOutstanding", mock.Anything, mock.Anything).Return(domain.ResponseLoanOutstanding{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Loan").Return(loanUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req = req.WithContext(candishared.SetToContext(req.Context(), candishared.ContextKeyTokenClaim, &candishared.TokenClaim{}))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.getLoanOutstanding(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}
