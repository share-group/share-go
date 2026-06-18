package secretutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(content string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(content))
	binaryData := h.Sum(nil)
	return hex.EncodeToString(binaryData)
}

func SHA256(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func Pkcs7AesEncrypt(text, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must be equal to AES block size")
	}

	plainText := pkcs7padding(text, aes.BlockSize)
	cipherText := make([]byte, len(plainText))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainText)

	// 将加密结果转换为 Base64 编码
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Pkcs7AesDecrypt(cipherTextBase64, key, iv []byte) (string, error) {
	if len(cipherTextBase64) <= 0 {
		return "", nil
	}

	// 解码 Base64 加密内容
	cipherText, err := base64.StdEncoding.DecodeString(string(cipherTextBase64))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must be equal to AES block size")
	}

	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	plainText := pkcs7UnPadding(cipherText)
	return string(plainText), nil
}

func pkcs7padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(plainText, padText...)
}

func pkcs7UnPadding(plainText []byte) []byte {
	length := len(plainText)
	unpadding := int(plainText[length-1])
	return plainText[:(length - unpadding)]
}

func ZeroAesEncrypt(plainText, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainText = ZeroPadding(plainText, block.BlockSize())
	cipherText := make([]byte, len(plainText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func ZeroAesDecrypt(cipherText, key, iv []byte) (string, error) {
	if len(cipherText) <= 0 {
		return "", nil
	}

	cipherText, err := base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(cipherText)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext length invalid")
	}
	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)
	plainText = ZeroUnPadding(plainText)
	return string(plainText), nil
}

func ZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	return append(data, padtext...)
}

func ZeroUnPadding(data []byte) []byte {
	i := len(data)
	for i > 0 && data[i-1] == 0 {
		i--
	}
	return data[:i]
}
