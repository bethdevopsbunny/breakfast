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

// Due to bcrypts default salting it wont work in the same way as other hashes
// this test function can tell you if a provided hash matches but the compute cannot be done beforehand
// sneaky stuff! try and find a way around this stuff.

func BCRYPTTEST(hash string, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
