package configs

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type PgClientCfg struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type DbConfig struct {
	UserTable    string
	MatrixTable  string
	TaskTable    string
	FinalTable   string
	SessionTable string
}

type AppConfig struct {
	Host       string
	Port       string
	ClaimsKey  string
	SignKey    string
	CookieName string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Config struct {
	PgClientCfg
	DbConfig
	AppConfig
}

func SetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	AccessTtl := time.Minute * 15
	RefreshTtl := time.Hour * 360
	return &Config{
		PgClientCfg: PgClientCfg{Host: os.Getenv("PGHOST"), Port: os.Getenv("PGPORT"),
			User: os.Getenv("PGUSER"), Password: os.Getenv("PGPASS"), Database: os.Getenv("PGDB")},
		AppConfig: AppConfig{Host: os.Getenv("HOST"), Port: os.Getenv("PORT"),
			ClaimsKey: os.Getenv("CLAIMSKEY"), SignKey: os.Getenv("SIGNKEY"), CookieName: os.Getenv("COOKIE"),
			AccessTTL: AccessTtl, RefreshTTL: RefreshTtl},
		DbConfig: DbConfig{UserTable: os.Getenv("USERSTAB"), MatrixTable: os.Getenv("MATRIXTAB"),
			TaskTable: os.Getenv("TASKTAB"), FinalTable: os.Getenv("FINALTAB"), SessionTable: os.Getenv("SESSIONTAB")},
	}
}
