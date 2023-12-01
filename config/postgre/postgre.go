package postgre

import (
	"fmt"
	"log"
	"time"

	"github.com/Mitra-Apps/user-service/domain/user/entity"
	"github.com/Mitra-Apps/user-service/lib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	username := lib.GetEnv("DB_USERNAME")
	password := lib.GetEnv("DB_PASSWORD")
	host := lib.GetEnv("DB_HOST")
	dbName := lib.GetEnv("DB_NAME")
	db, err := gorm.Open(postgres.Open("postgres://"+username+":"+password+"@"+host+"/"+dbName+"?sslmode=disable"),
		&gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalln(err)
	}

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
