package etc

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), 10)
}
