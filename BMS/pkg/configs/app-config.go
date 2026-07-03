package configs

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Connect() {
	conn, err := gorm.Open("mysql", "root:root charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	db = conn
}

// to use and call db connection in another file
func GetDB() *gorm.DB {
	return db
}
