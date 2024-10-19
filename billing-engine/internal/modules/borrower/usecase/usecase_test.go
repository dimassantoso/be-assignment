package usecase

import (
	"testing"

	mockdeps "github.com/golangid/candi/mocks/codebase/factory/dependency"
	"github.com/stretchr/testify/assert"
)

func TestNewBorrowerUsecase(t *testing.T) {
	mockDeps := &mockdeps.Dependency{}

	uc, setFunc := NewBorrowerUsecase(mockDeps)
	setFunc(nil)
	assert.NotNil(t, uc)
}
