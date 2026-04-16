package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	internalModel "github.com/coolpythoncodes/nigerian-universities/internal/model"
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

	if cfg.AppEnv == "production" {
		if err := validateSchema(db); err != nil {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
			return nil, err
		}
	} else {
		if err := db.AutoMigrate(&internalModel.Institution{}, &internalModel.User{}, &internalModel.ProductKey{}); err != nil {
			return nil, err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Verify connection works
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func validateSchema(db *gorm.DB) error {
	// Check if required tables/columns exist without modifying
	migrator := db.Migrator()

	// Ensure both core tables exist in production; fail fast if either is missing.
	if !migrator.HasTable(&internalModel.Institution{}) || !migrator.HasTable(&internalModel.User{}) || !migrator.HasTable(&internalModel.ProductKey{}) {
		return errors.New("required database tables are missing in production")
	}

	return nil
}

func getLogLevel(env string) logger.LogLevel {
	if env == "production" {
		return logger.Error // Only log errors in production
	}
	return logger.Info // Debug logging for dev/staging
}
