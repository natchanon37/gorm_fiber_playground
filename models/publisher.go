package models

import "gorm.io/gorm"

type Publisher struct {
	gorm.Model
	Details string
	Name    string
}

func CreatePublisher(db *gorm.DB, publisher *Publisher) error {
	result := db.Create(publisher)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
