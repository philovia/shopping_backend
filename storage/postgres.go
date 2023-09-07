package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		// "host=localhost port=5432 dbname=shopping user=postgres password=postgres sslmode=prefer connect_timeout=10")
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.DBName, config.User, config.Password, config.SSLMode)
	// fmt.Println(
	// 	DB_HOST=localhost,
	//     DB_PORT=5432,
	//     DB_USER=postgres,
	//     DB_PASS=postgres,
	//     DB_NAME=PostgreSQL15,
	//     DB_SSLMODE=disable,
	// )
	db, error := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if error != nil {
		return db, error
	}
	return db, nil
}
