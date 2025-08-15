package config

import (
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func Database() {
	// Connect to database
	Register(
		&db,
		&DBOption{
			Host:     confString("DB_HOST", "127.0.0.1"),
			Port:     confInt64("DB_PORT", int64(5432)),
			Database: confString("DB_DATABASE", "Database"),
			Username: confString("DB_USER", "postgres"),
			Password: confString("DB_PASSWORD", ""),
			TimeZone: confString("TZ", "Asia/Bangkok"),
			SSLMode:  confString("DB_SSLMODE", "disable"),
		},
	)
	logger.Info("database connected success")

}

type DBOption struct {
	DSN      string
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	TimeZone string
	SSLMode  string
}

func Register(conn **bun.DB, conf *DBOption) {
	if conf.DSN == "" {
		conf.DSN = generateDSN(conf)
	}

	config, err := pgx.ParseConfig(conf.DSN)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	sqldb := stdlib.OpenDB(*config)
	*conn = bun.NewDB(sqldb, pgdialect.New())

	if viper.GetBool("DEBUG") {
		(*conn).AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	err = (*conn).Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
}

func generateDSN(conf *DBOption) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
		conf.Host, conf.Port, conf.Database, conf.Username, conf.Password, conf.SSLMode, conf.TimeZone)
}

var (
	db     *bun.DB
	dbMap  = make(map[string]*bun.DB) // Initialize the dbMap
	dbLock sync.RWMutex
)

func GetDB() *bun.DB {
	return db
}

func DB(name ...string) *bun.DB {
	dbLock.RLock()
	defer dbLock.RUnlock()
	if dbMap == nil {
		panic("database not initialized") // Panic if dbMap is nil
	}
	if len(name) == 0 {
		return dbMap["default"] // Return the default database
	}

	db, ok := dbMap[name[0]]
	if !ok {
		panic("database not initialized") // Panic if the specified database is not found
	}
	return db
}

func Open(ctx context.Context) error {
	var err error
	for _, db := range dbMap {
		if errClose := db.Ping(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	return err
}

// Close function to close all registered databases
func Close(ctx context.Context) error {
	var err error
	for _, db := range dbMap {
		if errClose := db.Close(); errClose != nil {
			err = errors.Join(err, errClose)
		}
	}
	return err
}
