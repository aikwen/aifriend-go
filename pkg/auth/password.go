package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 将明文密码单向加密为 bcrypt 密文
func HashPassword(plaintextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 校验明文密码与数据库中的密文是否匹配
func CheckPassword(plaintextPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	// 密码完全正确
	return true, nil
}