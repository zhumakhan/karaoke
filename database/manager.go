package database

import (
	m "main/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Connect() (*gorm.DB, error) {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	database := os.Getenv("db_name")
	dialect := os.Getenv("db_type")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC", username, password, host, port, database)

	log.Println("Connecting to database.")
	log.Println(fmt.Sprintf("Database credentials: \n"+
		"	Host: %s\n"+
		"	Port: %s\n"+
		"	Database name: %s\n"+
		"	Username: %s\n"+
		"	Password: %s", host, port, database, username, password))
	var err error
	db, err = gorm.Open(dialect, dataSource)
	if err == nil {
		log.Println("Connected successfully.")
	} else {
		log.Println("Connection failed.")
	}
	db.LogMode(true)
	db.Set("gorm:table_options", "charset=utf8")
	//DropTables()
	CreateTables()
	CreateAlbums()
	db.Exec("ALTER TABLE musics ADD FULLTEXT (`title`);")
	return db, err
}

func CreateTables() {
	db.AutoMigrate(&m.Music{})
	db.AutoMigrate(&m.User{})
	db.AutoMigrate(&m.Admin{})
	db.AutoMigrate(&m.Author{})
	db.AutoMigrate(&m.Genre{})
	db.AutoMigrate(&m.Album{})
	db.AutoMigrate(&m.PlayListHelper{})
	db.AutoMigrate(&m.PlayList{})
	db.AutoMigrate(&m.UserRating{})
}

func DropTables() {
	db.DropTableIfExists(&m.Music{})
	db.DropTableIfExists(&m.User{})
	db.DropTableIfExists(&m.Author{})
	db.DropTableIfExists(&m.Genre{})
	db.DropTableIfExists(&m.Album{})
	db.DropTableIfExists(&m.PlayListHelper{})
	db.DropTableIfExists(&m.PlayList{})
	db.DropTableIfExists(&m.UserRating{})
}
func CreateAlbums(){
	db.Create(&m.Genre{Name : "Pop"})
	db.Create(&m.Genre{Name : "Indie"})
	db.Create(&m.Genre{Name : "Rock"})
	db.Create(&m.Genre{Name : "Electro"})
	db.Create(&m.Genre{Name : "Country"})
	db.Create(&m.Genre{Name : "Techno"})
	db.Create(&m.Genre{Name : "Dubstep"})
	db.Create(&m.Genre{Name : "Jazz"})

}