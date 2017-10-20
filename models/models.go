// Package models contains all structs used throughout the porgram
//
//
package models

// ResponseNewUUID is the struct of the response to /aws/new

type SignupData struct {
	Email    string `json:"email"`
	Password string `json:password`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:password`
}

type VerifyCaptchaResponse struct {
	Result bool `json:"success"`
}
