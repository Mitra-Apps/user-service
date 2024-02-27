package postgre

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBConn() (*gorm.DB, error) {
	if _, err := os.Stat("./../../../../.env"); !os.IsNotExist(err) {
		err := godotenv.Load(os.ExpandEnv("./../../../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
	}
	return LocalDatabase()
}

func LocalDatabase() (*gorm.DB, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME_TEST")
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
		log.Fatalf("failed to connect database: %v", err)
	}

	// db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	db.Migrator().DropTable("user_roles", &entity.Role{}, &entity.User{})
	db.AutoMigrate(&entity.User{}, &entity.Role{})

	return db, nil
}

func seedUser(db *gorm.DB) ([]*entity.User, error) {
	var users []*entity.User
	user := &entity.User{
		Name:     "name1",
		Email:    "test1@mail.com",
		Password: "password1",
	}
	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}
	users = append(users, user)
	return users, nil
}

func seedRole(db *gorm.DB) (*[]entity.Role, error) {
	role := []entity.Role{
		{
			RoleName: "customer",
		},
		{
			RoleName: "merchant",
		},
	}
	for _, v := range role {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return &role, nil
}
