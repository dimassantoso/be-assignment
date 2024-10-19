package repository

import (
	"billing-engine/pkg/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBillingRepoSQL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, _ := helper.MockGormDB()
		billingRepo := NewBillingRepoSQL(db, db)
		assert.NotNil(t, billingRepo)
	})
}
