package persistence

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	validator "github.com/go-playground/validator/v10"
)

type OpenConnectionInput struct {
	Config     *gorm.Config `json:"gormConfig" validation:"required"`
	DBHost     string       `json:"dbHost" validation:"required"`
	DBUsername string       `json:"dbUsername" validation:"required"`
	DBPassword string       `json:"dbPassword" validation:"required"`
	DBName     string       `json:"dbName" validation:"required"`
	DBPort     string       `json:"dbPort" validation:"required"`
}

func OpenConnection(input *OpenConnectionInput) (*gorm.DB, error) {
	connectionURI, err := newConnectionURI(input)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(
		postgres.Open(connectionURI),
		input.Config,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}

func newConnectionURI(input *OpenConnectionInput) (string, error) {
	err := validator.New().Struct(input)
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		input.DBHost,
		input.DBUsername,
		input.DBPassword,
		input.DBName,
		input.DBPort,
	)
	return uri, nil
}
