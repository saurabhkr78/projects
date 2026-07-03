package models

import (
	"BMS/pkg/configs"
	"github.com/jinzhu/gorm"
	"log"
)

type Book struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:100;not null" json:"name"`
	Author      string  `json:"author"`
	Publication string  `json:"publication"`
	Price       float64 `json:"price"`
}

var db *gorm.DB

/*
When the application starts, Connect() opens the database and stores the *gorm.DB object. GetDB() returns that same shared object. AutoMigrate() then compares the Book struct with the existing database schema and creates or updates the table if necessary. It does not recreate the table on every startup.

&Book{}: Create an empty Book value for the struct and give the pointer to that empty book
*/
func init() {
	configs.Connect()
	db = configs.GetDB() //points to db engine
	db.AutoMigrate(&Book{})
}

// create a new book in the database
func CreateBook(book *Book) *Book {
	db.NewRecord(book)
	db.Create(book)
	return book
}
func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books) //SELECT * FROM books;
	return Books
}
func DeleteBookByID(id int64) error {
	var book Book

	result := db.Where("id = ?", id).Delete(&book)
	if result.Error != nil {
		return result.Error
	}

	log.Println("Book deleted successfully")
	return nil
}
func FindBookByID(id uint) (*Book, error) {
	var book Book

	result := db.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &book, nil
}
func UpdateBookByID(id uint, book *Book) *Book {
	var existingBook Book

	result := db.First(&existingBook, id)
	if result.Error != nil {
		log.Println("Error finding book:", result.Error)
		return nil
	}

	existingBook.Name = book.Name
	existingBook.Author = book.Author
	existingBook.Publication = book.Publication
	existingBook.Price = book.Price

	result = db.Save(&existingBook)
	if result.Error != nil {
		log.Println("Error updating book:", result.Error)
		return nil
	}

	return &existingBook
}
