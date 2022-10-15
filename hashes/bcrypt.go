package hashes

import (
	"golang.org/x/crypto/bcrypt"
)

func BCRYPT(text string) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)

}
