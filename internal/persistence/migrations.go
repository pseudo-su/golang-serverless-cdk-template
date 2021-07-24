package persistence

import "gorm.io/gorm"

func MigrateAll(db *gorm.DB) error {

	err := db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`).Error

	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&League{},
	)

	if err != nil {
		return err
	}

	return nil
}
