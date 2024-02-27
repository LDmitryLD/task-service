package db

import (
	"database/sql"
	"fmt"
	"projects/LDmitryLD/task-service/task/config"
	"projects/LDmitryLD/task-service/task/internal/db/adapter"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewSqlDB(conf config.DB, logger *logrus.Logger) (*sqlx.DB, *adapter.SQLAdapter, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Password, conf.Name)

	var (
		dbRaw *sql.DB
		err   error
	)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second + time.Duration(conf.Timeout))

	for {
		select {
		case <-timeoutExceeded:
			return nil, nil, fmt.Errorf("db connection failed after %d timeout %s", conf.Timeout, err)
		case <-ticker.C:
			dbRaw, err = sql.Open(conf.Driver, dsn)
			if err != nil {
				return nil, nil, err
			}
			err = dbRaw.Ping()
			if err == nil {
				db := sqlx.NewDb(dbRaw, conf.Driver)
				db.SetMaxOpenConns(conf.MaxConn)
				db.SetMaxIdleConns(conf.MaxConn)
				sqlAdapter := adapter.NewSQLAdapter(db)
				return db, sqlAdapter, nil
			}
			logger.Error("failed to connect to the database:", err.Error(), "dsn:", dsn)
		}
	}
}
