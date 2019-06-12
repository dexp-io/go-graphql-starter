package dexp

import "golang.org/x/crypto/bcrypt"

func HashPassword(originalPassword string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(originalPassword), 10)
	return string(pass), err
}

func ComparePassword(password, hash string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil {
		return true
	}

	return false
}
