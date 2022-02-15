package db

import (
	"auth-service/utils"
	"context"
	"database/sql"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/go-sql-driver/mysql"
	"time"
)

func getConfig(cfg utils.Config) mysql.Config {
	config := mysql.Config{
		User:                 cfg.User,
		Passwd:               cfg.Passwd,
		Net:                  cfg.Net,
		Addr:                 cfg.Addr,
		DBName:               cfg.DBName,
		AllowNativePasswords: cfg.AllowNativePasswords,
	}
	return config
}

type BaseHandler struct {
	DB *sql.DB
}

func Connect(cfg utils.Config) (*sql.DB, error) {
	config := getConfig(cfg)
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	_, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)
	return db, nil
}

func Query(db *sql.DB, query string) (string, error) {
	var data string
	row := db.QueryRow(query)
	err := row.Scan(&data)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			err = fmt.Errorf("no rows were returned")
		default:
			err = fmt.Errorf("error: %s", err)
		}
	}
	return data, err
}

func VerifyPasswd(db *sql.DB, email string, pass string) (bool, error) {
	storedPass, err := Query(db, fmt.Sprintf("SELECT password FROM users WHERE email = '%s'", email))
	if err != nil {
		return false, err
	}
	match, err := argon2id.ComparePasswordAndHash(pass, storedPass)
	if err != nil {
		return false, err
	}
	return match, nil
}
