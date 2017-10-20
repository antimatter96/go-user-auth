package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"encoding/json"

	"../constants"
	"github.com/julienschmidt/httprouter"

	"../models"
	//"../db"
	//"strings"
	//"bytes"
	//"github.com/asaskevich/govalidator"
	//"crypto/rand"
	//"encoding/base64"
	//"time"

	"golang.org/x/crypto/bcrypt"
)

var signupTemplate *template.Template

func init() {
	signupTemplate = template.Must(template.ParseFiles("./template/home.html"))
}

func SignupHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	signupTemplate.Execute(w, nil)
}

func SignupHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var receivedData models.SignupData
	decoderReceivedData := json.NewDecoder(r.Body)
	errReceivedData := decoderReceivedData.Decode(&receivedData)
	defer r.Body.Close()
	if errReceivedData != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errReceivedData)
		return
	}

	passwordBytes := []byte(receivedData.Password)
	fmt.Println(passwordBytes)

	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, constants.BcryptCost)
	if err != nil {
		fmt.Println(err)
	}
	hashedString := string(hashedBytes)
	fmt.Println(hashedString)
	signupTemplate.Execute(w, nil)
}
