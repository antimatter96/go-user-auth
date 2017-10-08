package handlers

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"../constants"
	"encoding/json"
	"../models"
	"../db"
	"strings"
	"bytes"
	"github.com/asaskevich/govalidator"
	"crypto/rand"
	"encoding/base64"
	"time"
)

var CaptchaURL string
var CaptchaSecretKey string
var MissingFieldsError []byte
var EmailInvalidError []byte

var oneHour    time.Duration = 720 * time.Minute

func init(){
	CaptchaURL = constants.CaptchaURL
	CaptchaSecretKey = constants.CaptchaSecretKey
	MissingFieldsError = []byte("Missing Fields") 
	EmailInvalidError = []byte("Invalid Email") 
}

func CaptchaVerify(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var receivedData models.CaptchaData
	decoderReceivedData := json.NewDecoder(r.Body)
	errReceivedData := decoderReceivedData.Decode(&receivedData)
	defer r.Body.Close()
	if errReceivedData != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errReceivedData)
		return
	}
	
	if receivedData.CaptchaResponse == "" || receivedData.Username == "" || receivedData.Email == ""  || receivedData.Organisation == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Write(MissingFieldsError)
		fmt.Println("Missing Fields")
		return
	}
	
	if !govalidator.IsEmail(receivedData.Email){
		w.WriteHeader(http.StatusBadRequest)
		w.Write(EmailInvalidError)
		fmt.Println("Invalid Email", receivedData.Email)
		return
	}
	
	payload := strings.NewReader(getDataToSend(receivedData.CaptchaResponse))
	reqCaptchaVerify, _ := http.NewRequest("POST", CaptchaURL, payload)
	reqCaptchaVerify.Header.Add("content-type", "application/x-www-form-urlencoded")
	
	resCaptchaVerify, errResCaptcha := http.DefaultClient.Do(reqCaptchaVerify)
	if errResCaptcha!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(errResCaptcha)
		return
	}
	
	var verifyResponse models.VerifyCaptchaResponse
	decoderCaptchaVerify := json.NewDecoder(resCaptchaVerify.Body)
	errCaptchaVerify := decoderCaptchaVerify.Decode(&verifyResponse)
	defer resCaptchaVerify.Body.Close()
	if errCaptchaVerify != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, errCaptchaVerify)
		return
	}
	
	if verifyResponse.Result == true{
		newCSRF, _ := GenerateRandomString(32)
		plainValue := map[string]string{"value": newCSRF,}
		encodedValue, errCookieEncode := Cookie.Encode("_csrfToken", plainValue)
		if errCookieEncode != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//fmt.Println(errCookieEncode)
		} else{
			//fmt.Println(newCSRF,encodedValue)
			cookie := &http.Cookie{
				Name:  "_csrfToken",
				Value: encodedValue,
				Path:  "/",
				HttpOnly : true,
				Domain: "amazonaws.com",
				Expires: time.Now().Add(oneHour),
			} 
			http.SetCookie(w, cookie)
			w.Header().Set("X-Csrf-Token",newCSRF)
			w.WriteHeader(http.StatusOK)
		}
		go addTempUser(receivedData.Username,receivedData.Email,receivedData.Organisation)
		
	} else{
		w.WriteHeader(http.StatusForbidden)
	}

}

func getDataToSend(captchaResponse string) string {
	var buffer bytes.Buffer
	buffer.WriteString("secret=")
	buffer.WriteString(CaptchaSecretKey)
	buffer.WriteString("&response=")
	buffer.WriteString(captchaResponse)
	return buffer.String()
}

func addTempUser(name,email,org string){
	fmt.Println("Adding", email, name,org)
	done := db.AddTempData(name,email,org)
	if !done{
		fmt.Println("TEMP USER ERROR",name,email,org)
	}
}


func GenerateRandomString(length int) (string, error) {
	x := make([]byte, length)
	_, err := rand.Read(x)
   
	if err != nil {
		return "", err
	}
	
	return base64.URLEncoding.EncodeToString(x), err
}