package database

import (
	"belajar-golang-restful-api/helper"
	"database/sql"
	"fmt"
	"os"
	"time"
)

var (
	HOST     = os.Getenv("MYSQLHOST")
	USER     = os.Getenv("MYSQLUSER")
	PASSWORD = os.Getenv("MYSQLPASSWORD")
	DBPORT   = os.Getenv("MYSQLPORT")
	DBNAME   = os.Getenv("MYSQLDATABASE")
)

func NewDB() *sql.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USER, PASSWORD, HOST, DBPORT, DBNAME)
	db, err := sql.Open("mysql", config)
	helper.PanicError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
