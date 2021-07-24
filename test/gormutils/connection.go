package gormutils

import "gorm.io/gorm"

func DeleteDatabase(db *gorm.DB) error {
	// TODO: implement DeleteDatabase test Helper
	return nil
}

func CreateDatabase() error {
	// TODO: implement CreateDatabase test helper
	return nil
}

func DeleteAll(db *gorm.DB) error {
	tables := GetAllTableNames(db)
	ntables := len(tables)
	if ntables == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < ntables; i++ {
			tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("drop table " + tables[i] + " cascade")
		}
		return nil
	})
}

// ============================================================================

func GetAllTableNames(db *gorm.DB) []string {
	var tables []string
	if err := db.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		panic(err)
	}
	return tables
}
