package repository

import (
	"billing-engine/pkg/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBorrowerRepoSQL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, _ := helper.MockGormDB()
		borrowerRepo := NewBorrowerRepoSQL(db, db)
		assert.NotNil(t, borrowerRepo)
	})
}
