// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	context "context"

	candishared "github.com/golangid/candi/candishared"

	domain "billing-engine/internal/modules/billing/domain"

	mock "github.com/stretchr/testify/mock"

	shareddomain "billing-engine/pkg/shared/domain"
)

// BillingRepository is an autogenerated mock type for the BillingRepository type
type BillingRepository struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, filter
func (_m *BillingRepository) Find(ctx context.Context, filter *domain.FilterBilling) (shareddomain.Billing, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 shareddomain.Billing
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FilterBilling) (shareddomain.Billing, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FilterBilling) shareddomain.Billing); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(shareddomain.Billing)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.FilterBilling) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOverdueBilling provides a mock function with given fields: ctx, filter
func (_m *BillingRepository) FindOverdueBilling(ctx context.Context, filter *domain.FilterOverdueBilling) (shareddomain.OverdueBilling, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOverdueBilling")
	}

	var r0 shareddomain.OverdueBilling
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FilterOverdueBilling) (shareddomain.OverdueBilling, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FilterOverdueBilling) shareddomain.OverdueBilling); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(shareddomain.OverdueBilling)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.FilterOverdueBilling) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindPaymentMethod provides a mock function with given fields: ctx, filter
func (_m *BillingRepository) FindPaymentMethod(ctx context.Context, filter *domain.FilterPaymentMethod) (shareddomain.PaymentMethod, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindPaymentMethod")
	}

	var r0 shareddomain.PaymentMethod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FilterPaymentMethod) (shareddomain.PaymentMethod, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FilterPaymentMethod) shareddomain.PaymentMethod); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(shareddomain.PaymentMethod)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.FilterPaymentMethod) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, data, updateOptions
func (_m *BillingRepository) Save(ctx context.Context, data *shareddomain.Billing, updateOptions ...candishared.DBUpdateOptionFunc) error {
	_va := make([]interface{}, len(updateOptions))
	for _i := range updateOptions {
		_va[_i] = updateOptions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, data)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *shareddomain.Billing, ...candishared.DBUpdateOptionFunc) error); ok {
		r0 = rf(ctx, data, updateOptions...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveMany provides a mock function with given fields: ctx, data
func (_m *BillingRepository) SaveMany(ctx context.Context, data []shareddomain.Billing) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for SaveMany")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []shareddomain.Billing) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBillingRepository creates a new instance of BillingRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBillingRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BillingRepository {
	mock := &BillingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
