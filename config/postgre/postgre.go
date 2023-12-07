package postgre

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	conn := "postgres://" + username + ":" + password + "@" + host + "/" + dbName + "?sslmode=disable"
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalln(err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Tables has been migrated")

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour * 6)

	fmt.Println("Database successfully connected!")

	return db
}
