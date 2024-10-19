package repository

import (
	"billing-engine/pkg/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLoanRepoSQL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, _ := helper.MockGormDB()
		loanRepo := NewLoanRepoSQL(db, db)
		assert.NotNil(t, loanRepo)
	})
}
