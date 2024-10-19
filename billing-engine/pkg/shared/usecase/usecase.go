package usecase

import (
	"sync"

	// @candi:usecaseImport
	authusecase "billing-engine/internal/modules/auth/usecase"
	billingusecase "billing-engine/internal/modules/billing/usecase"
	borrowerusecase "billing-engine/internal/modules/borrower/usecase"
	loanusecase "billing-engine/internal/modules/loan/usecase"
	"billing-engine/pkg/shared/usecase/common"

	"github.com/golangid/candi/codebase/factory/dependency"
)

type (
	// Usecase unit of work for all usecase in modules
	Usecase interface {
		// @candi:usecaseMethod
		Auth() authusecase.AuthUsecase
		Billing() billingusecase.BillingUsecase
		Loan() loanusecase.LoanUsecase
		Borrower() borrowerusecase.BorrowerUsecase
	}

	usecaseUow struct {
		// @candi:usecaseField
		authusecase.AuthUsecase
		billingusecase.BillingUsecase
		loanusecase.LoanUsecase
		borrowerusecase.BorrowerUsecase
	}
)

var usecaseInst *usecaseUow
var once sync.Once

// SetSharedUsecase set singleton usecase unit of work instance
func SetSharedUsecase(deps dependency.Dependency) {
	once.Do(func() {
		usecaseInst = new(usecaseUow)
		var setSharedUsecaseFuncs []func(common.Usecase)
		var setSharedUsecaseFunc func(common.Usecase)

		// @candi:usecaseCommon
		usecaseInst.AuthUsecase, setSharedUsecaseFunc = authusecase.NewAuthUsecase(deps)
		setSharedUsecaseFuncs = append(setSharedUsecaseFuncs, setSharedUsecaseFunc)
		usecaseInst.BillingUsecase, setSharedUsecaseFunc = billingusecase.NewBillingUsecase(deps)
		setSharedUsecaseFuncs = append(setSharedUsecaseFuncs, setSharedUsecaseFunc)
		usecaseInst.LoanUsecase, setSharedUsecaseFunc = loanusecase.NewLoanUsecase(deps)
		setSharedUsecaseFuncs = append(setSharedUsecaseFuncs, setSharedUsecaseFunc)
		usecaseInst.BorrowerUsecase, setSharedUsecaseFunc = borrowerusecase.NewBorrowerUsecase(deps)
		setSharedUsecaseFuncs = append(setSharedUsecaseFuncs, setSharedUsecaseFunc)

		sharedUsecase := common.SetCommonUsecase(usecaseInst)
		for _, setFunc := range setSharedUsecaseFuncs {
			setFunc(sharedUsecase)
		}
	})
}

// GetSharedUsecase get usecase unit of work instance
func GetSharedUsecase() Usecase {
	return usecaseInst
}

// @candi:usecaseImplementation
func (uc *usecaseUow) Auth() authusecase.AuthUsecase {
	return uc.AuthUsecase
}

func (uc *usecaseUow) Billing() billingusecase.BillingUsecase {
	return uc.BillingUsecase
}

func (uc *usecaseUow) Loan() loanusecase.LoanUsecase {
	return uc.LoanUsecase
}

func (uc *usecaseUow) Borrower() borrowerusecase.BorrowerUsecase {
	return uc.BorrowerUsecase
}
