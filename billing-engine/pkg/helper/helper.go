package helper

import (
	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GeneratePassword(plainPassword string) string {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(hashPassword)
}

func MockGormDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mockDB, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	return gormDB, mockDB
}
