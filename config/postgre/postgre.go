package postgre

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connection() *gorm.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	portStr := ""
	if strings.Trim(port, " ") != "" {
		portStr = fmt.Sprintf("port=%s ", port)
	}
	dsn := fmt.Sprintf("host=%s "+portStr+"user=%s dbname=%s sslmode=disable password=%s", host, username, dbName, password)
	fmt.Printf("dsn: %s\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logrus.Panicf("failed to connect database: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// db.Migrator().DropTable("user_roles", &entity.Role{}, &entity.User{})
	err = db.AutoMigrate(&entity.User{}, &entity.Role{})
	if err != nil {
		logrus.Panicf("failed to migrate database: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		logrus.Panicf("failed to get database: %v", err)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour * 6)

	return db
}
