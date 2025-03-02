package utils

import (
	verifier "github.com/AfterShip/email-verifier"
)

func IsEmailValid(email string) bool {
	return verifier.IsAddressValid(email)
}

