package config

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	psql "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"strconv"
	"time"
)

var DB *gorm.DB

func getDBDNS(host, database, target string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s search_path=%s port=%s sslmode=%s TimeZone=%s target_session_attrs=%s",
		host,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		database,
		os.Getenv("DB_SCHEMA"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
		target)
}
func InitializeDatabase() {
	dns := getDBDNS(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_TARGET_SESSION"))
	replicaDSN := getDBDNS(os.Getenv("DB_REPLICA_HOST"), os.Getenv("DB_REPLICA_NAME"), os.Getenv("DB_REPLICA_TARGET_SESSION"))
	maxPoolSize, _ := strconv.Atoi(os.Getenv("DB_MAX_POOL_SIZE"))
	maxIdleSize, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_SIZE"))
	maxLeftTime, _ := strconv.Atoi(os.Getenv("DB_MAX_LEFT_TIME"))

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{postgres.Open(dns)},
			Replicas: []gorm.Dialector{postgres.Open(replicaDSN)},
			Policy:   dbresolver.RandomPolicy{},
		}).SetConnMaxLifetime(time.Hour * time.Duration(maxLeftTime)).
			SetMaxIdleConns(maxIdleSize).
			SetMaxOpenConns(maxPoolSize),
	)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	runMigrations(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(maxPoolSize)
	sqlDB.SetMaxIdleConns(maxIdleSize)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(maxLeftTime))

	DB = db
}

func runMigrations(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Driver()

	driver, err := psql.WithInstance(sqlDB, &psql.Config{})
	if err != nil {
		log.Printf("Migrations: %v", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		os.Getenv("DB_NAME"), driver,
	)
	if err != nil {
		log.Printf("Migrations: %v", err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrations: %v", err.Error())
	}
}
