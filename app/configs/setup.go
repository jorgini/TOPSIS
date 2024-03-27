package configs

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

type PostgresInstance struct {
	DB *sqlx.DB
}

func ConnectToDb(config *PgClientCfg) PostgresInstance {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		logrus.Fatalf("connecting to db failed with: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		logrus.Errorf("error occur on ping: %s", err.Error())
		err := db.Close()
		if err != nil {
			logrus.Fatal(err)
		}
		panic(err)
	}

	logrus.Info("Database connected successfully...")
	return PostgresInstance{DB: db}
}
