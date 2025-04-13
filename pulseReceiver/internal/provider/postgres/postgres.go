package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/lucasd-coder/pulseReceiver/config"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/logger"
	"github.com/lucasd-coder/pulseReceiver/internal/provider/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	initOnce sync.Once
	initErr  error
)

func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s search_path=%s",
		cfg.Host, cfg.PostgresPort, cfg.Username, cfg.Dbname, cfg.Password, cfg.Schema)
}

func setupConnectionPool(sqlDB *sql.DB, cfg *config.Config) {
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func StartDB(ctx context.Context, cfg *config.Config) error {
	initOnce.Do(func() {
		dsn := buildDSN(cfg)

		database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
		if err != nil {
			initErr = fmt.Errorf("could not connect to Postgres: %w", err)
			return
		}

		sqlDB, err := database.DB()
		if err != nil {
			initErr = fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
			return
		}

		setupConnectionPool(sqlDB, cfg)

		db = database
		logger.FromContext(ctx).Info("Postgres Connected")

		migrations.RunMigrations(db)
	})

	return initErr
}

func CloseConn() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func GetConn() *gorm.DB {
	return db
}
