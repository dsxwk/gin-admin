package pkg

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 对密码进行哈希加密
func BcryptHash(password string) string {
	if password == "" {
		return ""
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Md5 md5加密
func Md5(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

// Md5Salt md5加盐加密
func Md5Salt(str []byte, salt []byte) string {
	return Md5(append(str, salt...))
}
