package models

import "github.com/jinzhu/gorm"

// AutoMigrate creates the tables on DB
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(ItemStatus{}, Log{}, Item{})
}

// Seed calls the .Seed method on each model to populate DB
func Seed(db *gorm.DB, isProduction bool) {
	ItemStatus{}.Seed(db)
	if !isProduction {
		Item{}.Seed(db)
	}
}
