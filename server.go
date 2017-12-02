package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"./handlers"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
)

var fo *os.File
var homeTemplate *template.Template

func init() {
	homeTemplate = template.Must(template.ParseFiles("./template/home.html"))
	var errCreateFile error
	fo, errCreateFile = os.Create("output.txt")
	if errCreateFile != nil {
		fmt.Println(errCreateFile)
	}
	fmt.Println("Yo")
}

func main() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/login", handlers.LoginHandlerGet)
	router.POST("/login", handlers.LoginHandlerPost)
	router.GET("/signup", handlers.SignupHandlerGet)
	router.POST("/signup", handlers.SignupHandlerPost)

	router.ServeFiles("/static/*filepath", http.Dir("./template/static/"))

	loggedRouter := gorillaHandlers.LoggingHandler(fo, router)

	http.ListenAndServe(":8080", loggedRouter)
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Printf("%+v\n\n",r)
	homeTemplate.Execute(w, nil)
}
