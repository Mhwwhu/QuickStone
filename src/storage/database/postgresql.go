package database

import (
	"fmt"
	"time"

	"QuickStone/src/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Client *gorm.DB

func init() {
	var err error

	cfg := gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.EnvCfg.PostgreSQLTablePrefix,
		},
	}

	if Client, err = gorm.Open(
		postgres.Open(
			fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
				config.EnvCfg.PostgreSQLAddr,
				config.EnvCfg.PostgreSQLUser,
				config.EnvCfg.PostgreSQLPassword,
				config.EnvCfg.PostgreSQLDatabase,
				config.EnvCfg.PostgreSQLPort)),
		&cfg,
	); err != nil {
		panic(err)
	}

	sqlDB, err := Client.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Hour)
}
