package resthandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"billing-engine/internal/modules/borrower/domain"
	mockusecase "billing-engine/pkg/mocks/modules/borrower/usecase"
	mocksharedusecase "billing-engine/pkg/mocks/shared/usecase"

	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/candishared"
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
	mockMiddleware.On("HTTPPermissionACL", mock.Anything).Return(func(http.Handler) http.Handler { return nil })
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
	mockRoute.On("PUT", mock.Anything, mock.Anything, mock.Anything)
	mockRoute.On("DELETE", mock.Anything, mock.Anything, mock.Anything)
	handler.Mount(mockRoute)
}

func TestRestHandler_getAllBorrower(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", wantUsecaseError: nil, wantRespCode: http.StatusOK,
		},
		{
			name: "Testcase #2: Negative", reqBody: "?page=str", wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			borrowerUsecase := &mockusecase.BorrowerUsecase{}
			borrowerUsecase.On("GetAllBorrower", mock.Anything, mock.Anything).Return(domain.ResponseBorrowerList{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Borrower").Return(borrowerUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodGet, "/"+tt.reqBody, strings.NewReader(tt.reqBody))
			req = req.WithContext(candishared.SetToContext(req.Context(), candishared.ContextKeyTokenClaim, &candishared.TokenClaim{}))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.getAllBorrower(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_getDetailBorrowerByID(t *testing.T) {
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

			borrowerUsecase := &mockusecase.BorrowerUsecase{}
			borrowerUsecase.On("GetDetailBorrower", mock.Anything, mock.Anything).Return(domain.ResponseBorrower{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Borrower").Return(borrowerUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req = req.WithContext(candishared.SetToContext(req.Context(), candishared.ContextKeyTokenClaim, &candishared.TokenClaim{}))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.getDetailBorrowerByID(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_createBorrower(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", reqBody: `{"email": "test@test.com"}`, wantUsecaseError: nil, wantRespCode: http.StatusCreated,
		},
		{
			name: "Testcase #2: Negative", reqBody: `{"email": test@test.com}`, wantUsecaseError: nil, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", reqBody: `{"email": "test@test.com"}`, wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", reqBody: `{"email": "test@test.com"}`, wantUsecaseError: nil, wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			borrowerUsecase := &mockusecase.BorrowerUsecase{}
			borrowerUsecase.On("CreateBorrower", mock.Anything, mock.Anything).Return(domain.ResponseBorrower{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Borrower").Return(borrowerUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.createBorrower(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_updateBorrower(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", reqBody: `{"email": "test@test.com"}`, wantUsecaseError: nil, wantRespCode: http.StatusOK,
		},
		{
			name: "Testcase #2: Negative", reqBody: `{"email": test@test.com}`, wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", reqBody: `{"email": test@test.com}`, wantUsecaseError: nil, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", reqBody: `{"email": "test@test.com"}`, wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			borrowerUsecase := &mockusecase.BorrowerUsecase{}
			borrowerUsecase.On("UpdateBorrower", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Borrower").Return(borrowerUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req = req.WithContext(candishared.SetToContext(req.Context(), candishared.ContextKeyTokenClaim, &candishared.TokenClaim{}))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.updateBorrower(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_deleteBorrower(t *testing.T) {
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

			borrowerUsecase := &mockusecase.BorrowerUsecase{}
			borrowerUsecase.On("DeleteBorrower", mock.Anything, mock.Anything).Return(tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Borrower").Return(borrowerUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.deleteBorrower(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_delinquentCheck(t *testing.T) {
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

			borrowerUsecase := &mockusecase.BorrowerUsecase{}
			borrowerUsecase.On("DelinquentCheck", mock.Anything, mock.Anything).Return(domain.ResponseDelinquentCheck{}, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Borrower").Return(borrowerUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req = req.WithContext(candishared.SetToContext(req.Context(), candishared.ContextKeyTokenClaim, &candishared.TokenClaim{}))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.delinquentCheck(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}
