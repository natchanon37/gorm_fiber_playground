package models

import (
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	PublisherID uint
	Publisher   Publisher
	Authors     []Author `gorm:"many2many:author_books;"`
}

func CreateBookWithAuthor(db *gorm.DB, book *Book) error {
	result := db.Create(&book)
	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
		return result.Error
	}
	return nil
}

func GetBookByID(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.Preload("Authors").First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error getting book: %v", result.Error)
		return nil
	}
	return &book
}

func GetBooks(db *gorm.DB) []Book {
	var books []Book
	//use Preload to load the associated data, otherwise it will be lazy loaded which will return only attributes of the model book
	result := db.Preload("Authors").Find(&books)
	// result := db.Find(&books)
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
