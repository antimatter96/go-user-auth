// Package models contains all structs used throughout the porgram
//
//
package models

// ResponseNewUUID is the struct of the response to /aws/new

type VerifyCaptchaResponse struct {
	Result bool `json:"success"`
}