package main

import (
  "fmt"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "html/template"
)

var homeTemplate *template.Template
var loginTemplate *template.Template

func init(){
	homeTemplate = template.Must(template.ParseFiles("./template/home.html"))
	loginTemplate = template.Must(template.ParseFiles("./template/login.html"))
	fmt.Println("Yo")
}

func main() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/login", LoginHandlerGet)
	router.POST("/login", LoginHandlerPost)
	http.ListenAndServe(":8080", router)
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	homeTemplate.Execute(w, nil)
}

func LoginHandlerGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	loginTemplate.Execute(w, nil)
}

func LoginHandlerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()
	fmt.Println(r.Form)
	homeTemplate.Execute(w, nil)
}
