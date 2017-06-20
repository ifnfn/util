package system

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SecretPassword 对密码进行加密
func SecretPassword(secretKey, password string) string {
	sum := md5.Sum([]byte(secretKey + password))

	return hex.EncodeToString(sum[:])
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandomChString 生成指定长度的随机字母字符串
func GetRandomChString(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GetRandomNumberString 生成随机数字字符串
func GetRandomNumberString(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// HashString returns a hashed string and an error
func HashString(password string) (string, error) {
	key, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

// HashBytes returns a hashed byte array and an error
func HashBytes(password []byte) ([]byte, error) {
	key, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// MatchString returns true if the hash matches the password
func MatchString(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	return false
}

// MatchBytes returns true if the hash matches the password
func MatchBytes(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		return true
	}
	return false
}

// HmacSha1 returns
func HmacSha1(value, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
