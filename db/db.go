// Package db Contains all methods used by the other functions
package db

import (
	"database/sql"
	"fmt"
	"../constants"
	_ "github.com/go-sql-driver/mysql"
)

// The main db object
var db *sql.DB

// These are prepared statements which just need parameters
// Used to avoid multiple trips to db
var addTempUser *sql.Stmt

func init() {
	var err error
	db, err = sql.Open("mysql", constants.DBConnectionString)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(5)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	
	var errAddTempUser error
	addTempUser, errAddTempUser = db.Prepare("insert into `users` (`name`,`email`,`organisation`) values (?,?,?)")
	if errAddTempUser != nil {
		fmt.Println(errAddTempUser)
	}

}

func CheckStatus() bool {
	err := db.Ping()
	if err != nil {
		return false
	}
	return true
}

// AddTempData is called by CaptchaVerify to add the received name and email
func AddTempData(name, email,org string) bool {
	_, err := addTempUser.Exec(name, email,org)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}