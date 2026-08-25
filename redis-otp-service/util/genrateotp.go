package util

import (
	"github.com/google/uuid"
)

func GenerateOTP() string {
	//generate a random 6 digit number
	//but uuid genrates a random string of 32 characters so we will take the first 6 characters of the uuid string
	otp := uuid.New().String()[:6]
	return otp

}
