package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		// This enables GORM's built-in query logger so you can see
		// the actual SQL queries in your terminal — very helpful for debugging.
		Logger: logger.Default.LogMode(getLogLevel(cfg.AppEnv)),
		// PrepareStmt: true, // execute the SQL statement in the database and cache the statement
	})
	if err != nil {
		return nil, err
	}

	if err := validateSchema(db); err != nil {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Verify connection works
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return nil, err
	}

	return db, nil
}

func validateSchema(db *gorm.DB) error {
	// Check if required tables/columns exist without modifying
	if !db.Migrator().HasTable(&models.Universities{}) {
		return errors.New("universities table doesn't exist in production")
	}

	return nil
}

func getLogLevel(env string) logger.LogLevel {
	if env == "production" {
		return logger.Error // Only log errors in production
	}
	return logger.Info // Debug logging for dev/staging
}
