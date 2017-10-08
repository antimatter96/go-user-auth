package handlers

import (
	"github.com/gorilla/securecookie"
	"../constants"
)

var Cookie *securecookie.SecureCookie

func init(){
	var hashKey = []byte(constants.CookieHashKey)
	var blockKey = []byte(constants.CookieBlockKey)
	Cookie = securecookie.New(hashKey,blockKey)
	Cookie.MaxAge(43200)
}