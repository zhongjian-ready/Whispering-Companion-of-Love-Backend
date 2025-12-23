package database

import (
	"fmt"
	"miniapp-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// DSN format: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
