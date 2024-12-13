package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func encrypt(key []byte, aesIv []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext = pkcs5Padding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(plaintext))

	if len(aesIv) != block.BlockSize() {
		return nil, errors.New("IV must be the same length as the block size")
	}

	iv := make([]byte, aes.BlockSize)

	copy(iv, aesIv)

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func decrypt(key, aesIv []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(aesIv) != block.BlockSize() {
		return nil, errors.New("IV must be the same length as the block size")
	}
	iv := make([]byte, aes.BlockSize)

	copy(iv, aesIv)

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))

	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = pkcs5UnPadding(plaintext)

	return plaintext, nil
}

func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - (len(src) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

func EncryptWithBase64(key []byte, aesIv []byte, plaintext []byte) (string, error) {
	sign, err := encrypt(key, aesIv, plaintext)
	if err != nil {
		return "", err
	}
	base64Sign := base64.StdEncoding.EncodeToString(sign)
	return base64Sign, err

}

func DecryptWithBase64(key, aesIv []byte, ciphertext string) (string, error) {
	originCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	plaintext, err := decrypt(key, aesIv, originCiphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
