package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	db "../db"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	//"../constants"
	//"encoding/json"
	//"../models"
	//"../db"
	//"strings"
	//"bytes"
	//"github.com/asaskevich/govalidator"
	//"crypto/rand"
	//"encoding/base64"
	//"time"
)

var loginTemplate *template.Template

func init() {
	loginTemplate = template.Must(template.ParseFiles("./template/login.html"))
}

func LoginHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginTemplate.Execute(w, nil)
}

func LoginHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	errParseForm := r.ParseForm()
	if errParseForm != nil {
		fmt.Println(errParseForm)
		loginTemplate.Execute(w, constErrInternalError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || !govalidator.IsEmail(email) {
		loginTemplate.Execute(w, constErrEmailMissing)
		return
	}

	if password == "" {
		loginTemplate.Execute(w, constErrPasswordMissing)
		return
	}

	present, passwordHash, errGetPassword := getHashedPasswordFromDB(email)

	if !present {
		//fmt.Println("User not there")
		loginTemplate.Execute(w, constErrNotRegistered)
		return
	}

	if errGetPassword != nil {
		//fmt.Println("Internal")
		loginTemplate.Execute(w, constErrInternalError)
		return
	}

	verified := checkPassword(&password, &passwordHash)

	if !verified {
		//fmt.Println("Verified")
		loginTemplate.Execute(w, constErrPasswordMatchFailed)
		return
	}

	loginTemplate.Execute(w, nil)
}

func getHashedPasswordFromDB(email string) (bool, string, error) {
	password, err := db.GetPasswordHash(email)
	if err != nil {
		if password == "-" {
			return false, "", err
		}
		return true, "", err
	}
	return true, password, nil
}

func checkPassword(userPassword, savedPassword *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*savedPassword), []byte(*userPassword))
	if err != nil {
		return false
	}
	return true
}
