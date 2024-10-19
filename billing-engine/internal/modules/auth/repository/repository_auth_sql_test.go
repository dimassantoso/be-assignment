package repository

import (
	"billing-engine/internal/modules/auth/domain"
	"billing-engine/pkg/helper"
	shareddomain "billing-engine/pkg/shared/domain"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/candishared"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestNewAuthRepoSQL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, _ := helper.MockGormDB()
		authRepo := NewAuthRepoSQL(db, db)
		assert.NotNil(t, authRepo)
	})
}

func TestAuthRepoSQL_Find(t *testing.T) {
	t.Run("TestCase #1: Success", func(t *testing.T) {
		db, sql := helper.MockGormDB()
		var (
			request = domain.FilterAuth{
				ID:    candihelper.ToIntPtr(1),
				Email: "abc@example.com",
			}
			returnedRow = shareddomain.Auth{
				ID:        1,
				Email:     "abc@example.com",
				Password:  "hashed_password",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			expectedQuery = `SELECT (.+) FROM "` + shareddomain.Auth{}.TableName() + `" WHERE (.+)`
		)

		sql.ExpectQuery(expectedQuery).WillReturnRows(sql.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).AddRow(returnedRow.ID, returnedRow.Email, returnedRow.Password, returnedRow.CreatedAt, returnedRow.UpdatedAt))

		repo := &authRepoSQL{readDB: db, writeDB: db}
		res, err := repo.Find(context.Background(), &request)
		assert.Equal(t, returnedRow, res)
		assert.NoError(t, err)
		if errExpect := sql.ExpectationsWereMet(); errExpect != nil {
			t.Errorf("expectation and result doest not match : %s", errExpect)
		}
	})

	t.Run("TestCase #2: Failed", func(t *testing.T) {
		db, sql := helper.MockGormDB()
		var (
			ctx     = candishared.SetToContext(context.Background(), candishared.ContextKeySQLTransaction, db)
			request = domain.FilterAuth{
				ID:    candihelper.ToIntPtr(1),
				Email: "abc@example.com",
			}

			expectedQuery = `SELECT (.+) FROM "` + shareddomain.Auth{}.TableName() + `" WHERE (.+)`
		)

		sql.ExpectQuery(expectedQuery).WillReturnError(gorm.ErrInvalidData)

		repo := &authRepoSQL{readDB: db, writeDB: db}
		_, err := repo.Find(ctx, &request)
		assert.Error(t, err)
		if errExpect := sql.ExpectationsWereMet(); errExpect != nil {
			t.Errorf("expectation and result doest not match : %s", errExpect)
		}
	})
}

func TestAuthRepoSQL_Save(t *testing.T) {
	t.Run("TestCase #1: Success (Insert)", func(t *testing.T) {
		db, sql := helper.MockGormDB()
		var (
			ctx  = candishared.SetToContext(context.Background(), candishared.ContextKeySQLTransaction, db)
			data = shareddomain.Auth{
				Email: "me@example.com",
			}
			expectedQuery = `INSERT INTO "` + shareddomain.Auth{}.TableName() + `" (.+) RETURNING "id"`
		)

		sql.ExpectBegin()
		sql.ExpectQuery(expectedQuery).WithArgs(data.Email, "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sql.NewRows([]string{"id"}).AddRow(1))
		sql.ExpectCommit()

		repo := &authRepoSQL{readDB: db, writeDB: db}
		err := repo.Save(ctx, &data)
		assert.NoError(t, err)

		if errExpect := sql.ExpectationsWereMet(); errExpect != nil {
			t.Errorf("expectation and result does not match: %s", errExpect)
		}
	})

	t.Run("TestCase #2: Failed", func(t *testing.T) {
		db, sql := helper.MockGormDB()
		var (
			ctx  = candishared.SetToContext(context.Background(), candishared.ContextKeySQLTransaction, db)
			data = shareddomain.Auth{
				Email: "me@example.com",
			}
			expectedQuery = `INSERT INTO "` + shareddomain.Auth{}.TableName() + `" (.+) RETURNING "id"`
		)

		sql.ExpectBegin()
		sql.ExpectQuery(expectedQuery).WithArgs(data.Email, "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(gorm.ErrInvalidDB)
		sql.ExpectRollback()

		repo := &authRepoSQL{readDB: db, writeDB: db}
		err := repo.Save(ctx, &data)
		assert.Error(t, err)

		if errExpect := sql.ExpectationsWereMet(); errExpect != nil {
			t.Errorf("expectation and result does not match: %s", errExpect)
		}
	})
}
