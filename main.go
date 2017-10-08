package main

import (
	"fmt"
	"net/http"
	
	"./constants"
	"./handlers"
	"github.com/julienschmidt/httprouter"
	gorillaHandlers "github.com/gorilla/handlers"
	"os"
	//"regexp"
)

var fo *os.File 
var allowedHosts map[string]bool

func init(){
	var errCreateFile error
	fo, errCreateFile = os.Create("output.txt")
	if errCreateFile!=nil{
		fmt.Println(errCreateFile)
	}
	allowedHosts = make(map[string]bool)
	for _,host := range constants.AllowedHosts{
		allowedHosts[host] = true
	}
}

func DisallowOthers(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		headerValue := r.Header.Get("X-Csrf-Token")
		if headerValue == "" || len(headerValue) < 30 {
			fmt.Println("NO HEADER OR Header FORMAT WROng");
			http.Error(w, UnauthorizedStatusText, http.StatusUnauthorized)
		} else{
			cookie, err := r.Cookie("_csrfToken")
			if err != nil{
				fmt.Println("NO COOKIE",err);
				http.Error(w, UnauthorizedStatusText, http.StatusUnauthorized)
			} else{
				value := make(map[string]string)
				err = handlers.Cookie.Decode("_csrfToken", cookie.Value, &value)
				if err != nil {
					fmt.Println("DECODE ERROR",err);
					http.Error(w, UnauthorizedStatusText, http.StatusUnauthorized)
				} else{
					if value["value"] == headerValue{
						h(w, r, ps)
					} else{
						fmt.Println("Header does not match cookie" , "Intrusion")
						http.Error(w, UnauthorizedStatusText, http.StatusUnauthorized)
					}
				}
			}
		}
	}
}

func BasicSecurity(handler http.Handler) http.Handler {
    AddHeaders := func(w http.ResponseWriter, r *http.Request) {
		if allowedHosts[r.Header.Get("origin")]{
			w.Header().Add("Access-Control-Allow-Origin",r.Header.Get("origin"))
		}
		w.Header().Add("Access-Control-Allow-Credentials","true")
		w.Header().Add("Access-Control-Allow-Headers","Content-Type,X-Csrf-Token")
		w.Header().Add("Access-Control-Expose-Headers","X-Csrf-Token")
	
		w.Header().Add("x-xss-protection", "1; mode=block")
		w.Header().Add("x-frame-options", "SAMEORIGIN")
		w.Header().Add("x-content-type", "nosniff")
        handler.ServeHTTP(w, r)
    }
    return http.HandlerFunc(AddHeaders)
}

var UnauthorizedStatusText string

func init(){
	UnauthorizedStatusText = http.StatusText(http.StatusUnauthorized)
}

func main() {

	router := httprouter.New()

	router.GET("/", handlers.HomeHandler) // This will be removed as front-end is static page
	router.POST("/captcha-verify", handlers.CaptchaVerify)
	router.GET("/aws/healthcheck", handlers.HealthCheck)
	
	// All Beyond this point require auth/CSRF/CaptchaVerify
	// 
	
	router.GET("/aws/new", DisallowOthers(handlers.AWSNew))
	router.GET("/aws/account/:UUID", DisallowOthers(handlers.AWSCheck))

	router.POST("/aws/account/register", DisallowOthers(handlers.AWSRegister))
	
	// DEV ONLY
	router.ServeFiles("/static/*filepath", http.Dir("./template/static/"))

	
	loggedRouter := gorillaHandlers.LoggingHandler(fo, router)
	loggedAndSecuredRouter := BasicSecurity(loggedRouter)
	
	fmt.Println(http.ListenAndServe(":" + constants.Port, loggedAndSecuredRouter))
}