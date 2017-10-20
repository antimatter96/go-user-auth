package handlers

import (
	"fmt"
	"html/template"
	"net/http"

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
	loginTemplate = template.Must(template.ParseFiles("./template/home.html"))
}

func LoginHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginTemplate.Execute(w, nil)
}

func LoginHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	fmt.Println(r.Form)
	loginTemplate.Execute(w, nil)
}
