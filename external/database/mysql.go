package database

import (
	"fmt"
	"time"

	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/jmoiron/sqlx"
	_"github.com/go-sql-driver/mysql"
)

//2
func ConnectSQL(cfgYAML configs.DBConfig) (db *sqlx.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfgYAML.User,
		cfgYAML.Password,
		cfgYAML.Host,
		cfgYAML.Port,
		cfgYAML.Name,
	)
	// fmt.Println("DSN:", dsn) // debug

	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	//cek konesksi
	if err = db.Ping(); err != nil {
		return nil, err
	}

	//dbpool
	db.SetMaxOpenConns(int(cfgYAML.DBConnPoll.MaxOpenConnection))
	db.SetMaxIdleConns(int(cfgYAML.DBConnPoll.MaxIdleConnection))
	db.SetConnMaxIdleTime(time.Duration(cfgYAML.DBConnPoll.MaxIdleTimeConnection) * time.Second)
	db.SetConnMaxLifetime(time.Duration(cfgYAML.DBConnPoll.MaxLifeTimeConnection) + time.Second)

	return db, nil
}
