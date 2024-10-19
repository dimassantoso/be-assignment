package resthandler

import (
	"billing-engine/internal/modules/auth/domain"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockusecase "billing-engine/pkg/mocks/modules/auth/usecase"
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
	responseLogin                       domain.ResponseLogin
	attemptLogin                        int
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

func TestRestHandler_createAuth(t *testing.T) {
	tests := []testCase{
		{
			name: "Testcase #1: Positive", reqBody: `{"email": "test@test.com", "password": "123!"}`, wantUsecaseError: nil, wantRespCode: http.StatusCreated,
		},
		{
			name: "Testcase #2: Negative", reqBody: `{"email": test@test.com, "password": "123!"}`, wantUsecaseError: nil, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #3: Negative", reqBody: `{"email": "test@test.com", "password": "123!"}`, wantUsecaseError: errFoo, wantRespCode: http.StatusBadRequest,
		},
		{
			name: "Testcase #4: Negative", reqBody: `{"email": "test@test.com", "password": "123!"}`, wantUsecaseError: nil, wantValidateError: errFoo, wantRespCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			authUsecase := &mockusecase.AuthUsecase{}
			authUsecase.On("CreateAuth", mock.Anything, mock.Anything).Return(tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Auth").Return(authUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.createAuth(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}

func TestRestHandler_login(t *testing.T) {
	tests := []testCase{
		{
			name:             "Testcase #1: Positive",
			reqBody:          `{"email": "test@test.com", "password": "123!", "keep_sign_in": true}`,
			wantUsecaseError: nil,
			responseLogin: domain.ResponseLogin{
				ID:    1,
				Email: "abc@example.com",
				Token: "generated_token",
			},
			attemptLogin: 5,
			wantRespCode: http.StatusOK,
		},
		{
			name:             "Testcase #2: Negative",
			reqBody:          `{"email": test@test.com, "password": "123!", "keep_sign_in": true}`,
			wantUsecaseError: nil,
			wantRespCode:     http.StatusBadRequest,
		},
		{
			name:             "Testcase #3: Negative",
			reqBody:          `{"email": "test@test.com", "password": "123!", "keep_sign_in": true}`,
			wantUsecaseError: errFoo,
			attemptLogin:     0,
			wantRespCode:     http.StatusBadRequest,
		},
		{
			name:              "Testcase #4: Negative",
			reqBody:           `{"email": "test@test.com", "password": "123!", "keep_sign_in": true}`,
			wantUsecaseError:  nil,
			wantValidateError: errFoo,
			attemptLogin:      0,
			wantRespCode:      http.StatusBadRequest,
		},
		{
			name:             "Testcase #5: Negative",
			reqBody:          `{"email": "test@test.com", "password": "123!", "keep_sign_in": true}`,
			wantUsecaseError: errors.New("password not match"),
			responseLogin:    domain.ResponseLogin{},
			attemptLogin:     4,
			wantRespCode:     http.StatusForbidden,
		},
		{
			name:             "Testcase #6: Negative",
			reqBody:          `{"email": "test@test.com", "password": "123!", "keep_sign_in": true}`,
			wantUsecaseError: errors.New("user not found"),
			responseLogin:    domain.ResponseLogin{},
			attemptLogin:     0,
			wantRespCode:     http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			authUsecase := &mockusecase.AuthUsecase{}
			authUsecase.On("LoginAuth", mock.Anything, mock.Anything).Return(domain.ResponseLogin{}, tt.attemptLogin, tt.wantUsecaseError)
			mockValidator := &mockinterfaces.Validator{}
			mockValidator.On("ValidateDocument", mock.Anything, mock.Anything).Return(tt.wantValidateError)

			uc := &mocksharedusecase.Usecase{}
			uc.On("Auth").Return(authUsecase)

			handler := RestHandler{uc: uc, validator: mockValidator}

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.reqBody))
			req.Header.Add(candihelper.HeaderContentType, candihelper.HeaderMIMEApplicationJSON)
			res := httptest.NewRecorder()
			handler.login(res, req)
			assert.Equal(t, tt.wantRespCode, res.Code)
		})
	}
}
