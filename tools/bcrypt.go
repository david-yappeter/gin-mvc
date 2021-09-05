package tools

import "golang.org/x/crypto/bcrypt"

func HashPassword(s string) string {
	raw, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(raw)
}

func ComparePass(hashed string, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
}
