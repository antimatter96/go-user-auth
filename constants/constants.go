// Pacakge contants stores all the constants used all over the service

package constants

import (
	"fmt"
	//"os"
	"github.com/spf13/viper"
)

/*
var ApiLink string
var RedisAddress string
var CaptchaSecretKey string
var CaptchaURL string

var CookieHashKey string
var CookieBlockKey string
var Port string
var AllowedHosts []string
*/

var DBConnectionString string
var BcryptCost int

func init() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	errViperRead := viper.ReadInConfig()
	if errViperRead != nil {
		panic(fmt.Errorf("fatal error config file: %s", errViperRead))
	}

	DBConnectionString = viper.GetString("db-connection-string")
	BcryptCost = viper.GetInt("bcrypt-cost")

	/*
		ApiLink = viper.GetString("actualApiLink")
		RedisAddress = viper.GetString("redisAddress")
		CaptchaSecretKey = viper.GetString("captcha.secret-key")
		CaptchaURL = viper.GetString("captcha.url")

		CookieHashKey = viper.GetString("cookie.hashKey")
		CookieBlockKey = viper.GetString("cookie.blockKey")
		Port = viper.GetString("port")
		AllowedHosts = viper.GetStringSlice("allowed-domains")
	*/
}
