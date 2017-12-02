package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"../constants"
	"github.com/julienschmidt/httprouter"

	db "../db"
	//"strings"
	//"bytes"
	"github.com/asaskevich/govalidator"
	//"crypto/rand"
	//"encoding/base64"
	//"time"

	"golang.org/x/crypto/bcrypt"
)

var signupTemplate *template.Template

func init() {
	signupTemplate = template.Must(template.ParseFiles("./template/signup.html"))
}

// S
func SignupHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	signupTemplate.Execute(w, nil)
}

// S
func SignupHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	errParseForm := r.ParseForm()
	if errParseForm != nil {
		fmt.Println(errParseForm)
		signupTemplate.Execute(w, constErrInternalError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || !govalidator.IsEmail(email) {
		signupTemplate.Execute(w, constErrEmailMissing)
		return
	}

	userPresent, errPresent := checkUser(email)
	if errPresent != nil {
		signupTemplate.Execute(w, constErrInternalError)
		return
	}
	if userPresent {
		signupTemplate.Execute(w, constErrEmailTaken)
		return
	}

	if password == "" {
		signupTemplate.Execute(w, constErrPasswordMissing)
		return
	}

	hashedString, errBcrypt := getHashedPassword(password)
	if errBcrypt != nil {
		signupTemplate.Execute(w, constErrInternalError)
		return
	}
	fmt.Println(hashedString)
	userAdded := addUser(email, hashedString)
	if !userAdded {
		signupTemplate.Execute(w, constErrInternalError)
		return
	}
	http.Redirect(w, r, "./login?success=true", http.StatusSeeOther)
}

func getHashedPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedBytes, errBcrypt := bcrypt.GenerateFromPassword(passwordBytes, constants.BcryptCost)
	if errBcrypt != nil {
		fmt.Println(errBcrypt)
		return "", errBcrypt
	}
	return string(hashedBytes), nil
}

func checkUser(email string) (bool, error) {
	present, err := db.CheckUser(email)
	if err != nil {
		return true, err
	}
	return present, nil
}

func addUser(email, password string) bool {
	return db.AddUser(email, password)
}
