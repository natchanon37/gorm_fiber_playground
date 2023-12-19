package models

import (
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string
	Author      string
	Description string
	Price       uint
}

func CreateBook(db *gorm.DB, book *Book) error {
	result := db.Create(&book)
	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
		return result.Error
	}
	return nil
}

func GetBookByID(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error getting book: %v", result.Error)
		return nil
	}
	return &book
}

func GetBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error getting books: %v", result.Error)
		return nil
	}
	return books
}

func UpdateBookByID(db *gorm.DB, book *Book) error {
	result := db.Model(&book).Updates(book)
	if result.Error != nil {
		log.Fatalf("Error updating book: %v", result.Error)
		return result.Error
	}
	return nil
}

func DeleteBookbyID(db *gorm.DB, id uint) error {
	result := db.Delete(&Book{}, id)
	if result.Error != nil {
		log.Fatalf("Error deleting book: %v", result.Error)
		return result.Error
	}
	return nil
}
