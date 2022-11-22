package password

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// CompareHashAndPassword ..
func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return true, nil
}
