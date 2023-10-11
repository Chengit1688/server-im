package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func Md5(s string, salt ...string) string {
	h := md5.New()
	h.Write([]byte(s))
	if len(salt) > 0 {
		h.Write([]byte(salt[0]))
	}
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}

func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("encrypt error")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func GetPassword(password, salt string) (pass string) {
	newPasswordFirst := password + salt
	passwordData := []byte(newPasswordFirst)
	has := md5.Sum(passwordData)
	pass = fmt.Sprintf("%x", has)
	return
}

func CheckPassword(dbpassword, userpassword, salt string) bool {
	newPasswordFirst := userpassword + salt
	passwordData := []byte(newPasswordFirst)
	has := md5.Sum(passwordData)
	password := fmt.Sprintf("%x", has)
	return dbpassword == password
}

func Encrypt(data []byte, key []byte) (content string, err error) {
	encrypt, err := AesEncrypt(data, key)
	if err != nil {
		return
	}

	content = base64.StdEncoding.EncodeToString(encrypt)
	return
}

func Decrypt(dataStr string, key []byte) (content string, err error) {
	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return
	}

	decrypt, err := AesDecrypt(data, key)
	if err != nil {
		return
	}
	content = string(decrypt)
	return
}
