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
var addUser *sql.Stmt
var getPassword *sql.Stmt
var checkUserStmt *sql.Stmt

func init() {
	var err error
	db, err = sql.Open("mysql", constants.DBConnectionString)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(3)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	var errAddUser error
	addUser, errAddUser = db.Prepare("insert into `users` (`email`,`password`) values (?,?)")
	if errAddUser != nil {
		fmt.Println(errAddUser)
	}

	var errGetPassword error
	getPassword, errGetPassword = db.Prepare("select `password` from `users` where `email` = ?")
	if errGetPassword != nil {
		fmt.Println(errGetPassword)
	}

	var errCheckUser error
	checkUserStmt, errCheckUser = db.Prepare("select `id` from `users` where `email` = ?")
	if errCheckUser != nil {
		fmt.Println(errCheckUser)
	}

}

func CheckStatus() bool {
	err := db.Ping()
	if err != nil {
		return false
	}
	return true
}

func AddUser(email, password string) bool {
	_, err := addUser.Exec(email, password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetPasswordHash(email string) (string, error) {
	var password string
	err := getPassword.QueryRow(email).Scan(&password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "-", err
		}
		return "", err
	}
	return password, nil
}

func CheckUser(email string) (bool, error) {
	var userId int
	err := checkUserStmt.QueryRow(email).Scan(&userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		fmt.Println(err)
		return false, err

	}
	if userId == 0 {
		return false, nil
	}
	return true, nil
}
