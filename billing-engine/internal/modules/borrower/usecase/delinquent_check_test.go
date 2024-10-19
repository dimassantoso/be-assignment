package usecase

import (
	mockbillingrepo "billing-engine/pkg/mocks/modules/billing/repository"
	mockrepo "billing-engine/pkg/mocks/modules/borrower/repository"
	mocksharedrepo "billing-engine/pkg/mocks/shared/repository"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestBorrowerUsecaseImpl_DelinquentCheck(t *testing.T) {
	t.Run("Testcase #1: Positive", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{ID: 10}, nil)

		billingRepo := &mockbillingrepo.BillingRepository{}
		billingRepo.On("FindOverdueBilling", mock.Anything, mock.Anything).Return(shareddomain.OverdueBilling{
			ID:           1,
			Name:         "Example",
			CountOverdue: 3,
		}, nil)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)
		repoSQL.On("BillingRepo").Return(billingRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.DelinquentCheck(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("Testcase #2: Positive", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{ID: 10}, nil)

		billingRepo := &mockbillingrepo.BillingRepository{}
		billingRepo.On("FindOverdueBilling", mock.Anything, mock.Anything).Return(shareddomain.OverdueBilling{}, gorm.ErrRecordNotFound)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)
		repoSQL.On("BillingRepo").Return(billingRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.DelinquentCheck(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("Testcase #3: Negative", func(t *testing.T) {

		borrowerRepo := &mockrepo.BorrowerRepository{}
		borrowerRepo.On("Find", mock.Anything, mock.Anything).Return(shareddomain.Borrower{}, gorm.ErrRecordNotFound)

		billingRepo := &mockbillingrepo.BillingRepository{}
		billingRepo.On("FindOverdueBilling", mock.Anything, mock.Anything).Return(shareddomain.OverdueBilling{}, gorm.ErrRecordNotFound)

		repoSQL := &mocksharedrepo.RepoSQL{}
		repoSQL.On("BorrowerRepo").Return(borrowerRepo)
		repoSQL.On("BillingRepo").Return(billingRepo)

		uc := borrowerUsecaseImpl{
			repoSQL: repoSQL,
		}

		_, err := uc.DelinquentCheck(context.Background(), 1)
		assert.Error(t, err)
	})
}
