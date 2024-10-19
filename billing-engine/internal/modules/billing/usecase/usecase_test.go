package usecase

import (
	"testing"

	mockdeps "github.com/golangid/candi/mocks/codebase/factory/dependency"
	"github.com/stretchr/testify/assert"
)

func TestNewBillingUsecase(t *testing.T) {
	mockDeps := &mockdeps.Dependency{}

	uc, setFunc := NewBillingUsecase(mockDeps)
	setFunc(nil)
	assert.NotNil(t, uc)
}
