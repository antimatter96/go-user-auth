// Pacakge contants stores all the constants used all over the service

package constants

import (
	"fmt"
	//"os"
	"github.com/spf13/viper"
)


var ApiLink string
var RedisAddress string
var CaptchaSecretKey string
var CaptchaURL string
var DBConnectionString string
var CookieHashKey string
var CookieBlockKey string
var Port string
var AllowedHosts []string

func init() {
	
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	errViperRead := viper.ReadInConfig()
	if errViperRead != nil {
		panic(fmt.Errorf("fatal error config file: %s", errViperRead))
	}

	ApiLink = viper.GetString("actualApiLink")
	RedisAddress = viper.GetString("redisAddress")
	CaptchaSecretKey = viper.GetString("captcha.secret-key")
	CaptchaURL = viper.GetString("captcha.url")
	DBConnectionString = viper.GetString("db-connection-string")
	CookieHashKey = viper.GetString("cookie.hashKey")
	CookieBlockKey = viper.GetString("cookie.blockKey")
	Port = viper.GetString("port")
	AllowedHosts = viper.GetStringSlice("allowed-domains")
}
