package usecase

import (
	"testing"

	mockdeps "github.com/golangid/candi/mocks/codebase/factory/dependency"
	mockinterfaces "github.com/golangid/candi/mocks/codebase/interfaces"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthUsecase(t *testing.T) {
	mockPublisher := &mockinterfaces.Publisher{}
	mockBroker := &mockinterfaces.Broker{}
	mockBroker.On("GetPublisher").Return(mockPublisher)

	mockCache := &mockinterfaces.Cache{}
	mockRedisPool := &mockinterfaces.RedisPool{}
	mockRedisPool.On("Cache").Return(mockCache)

	mockDeps := &mockdeps.Dependency{}
	mockDeps.On("GetRedisPool").Return(mockRedisPool)

	uc, setFunc := NewAuthUsecase(mockDeps)
	setFunc(nil)
	assert.NotNil(t, uc)
}
