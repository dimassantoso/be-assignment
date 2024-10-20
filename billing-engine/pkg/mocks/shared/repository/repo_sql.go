// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	billingrepository "billing-engine/internal/modules/billing/repository"
	borrowerrepository "billing-engine/internal/modules/borrower/repository"

	context "context"

	loanrepository "billing-engine/internal/modules/loan/repository"

	mock "github.com/stretchr/testify/mock"

	repository "billing-engine/internal/modules/auth/repository"
)

// RepoSQL is an autogenerated mock type for the RepoSQL type
type RepoSQL struct {
	mock.Mock
}

// AuthRepo provides a mock function with given fields:
func (_m *RepoSQL) AuthRepo() repository.AuthRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AuthRepo")
	}

	var r0 repository.AuthRepository
	if rf, ok := ret.Get(0).(func() repository.AuthRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.AuthRepository)
		}
	}

	return r0
}

// BillingRepo provides a mock function with given fields:
func (_m *RepoSQL) BillingRepo() billingrepository.BillingRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BillingRepo")
	}

	var r0 billingrepository.BillingRepository
	if rf, ok := ret.Get(0).(func() billingrepository.BillingRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(billingrepository.BillingRepository)
		}
	}

	return r0
}

// BorrowerRepo provides a mock function with given fields:
func (_m *RepoSQL) BorrowerRepo() borrowerrepository.BorrowerRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BorrowerRepo")
	}

	var r0 borrowerrepository.BorrowerRepository
	if rf, ok := ret.Get(0).(func() borrowerrepository.BorrowerRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(borrowerrepository.BorrowerRepository)
		}
	}

	return r0
}

// LoanRepo provides a mock function with given fields:
func (_m *RepoSQL) LoanRepo() loanrepository.LoanRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LoanRepo")
	}

	var r0 loanrepository.LoanRepository
	if rf, ok := ret.Get(0).(func() loanrepository.LoanRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(loanrepository.LoanRepository)
		}
	}

	return r0
}

// WithTransaction provides a mock function with given fields: ctx, txFunc
func (_m *RepoSQL) WithTransaction(ctx context.Context, txFunc func(context.Context) error) error {
	ret := _m.Called(ctx, txFunc)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, txFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepoSQL creates a new instance of RepoSQL. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepoSQL(t interface {
	mock.TestingT
	Cleanup(func())
}) *RepoSQL {
	mock := &RepoSQL{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
