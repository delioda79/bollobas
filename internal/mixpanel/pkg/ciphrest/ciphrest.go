package ciphrest

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var key []byte
var iv []byte
var block cipher.Block

//InitCipher sets key and init vector
func InitCipher(rawKey string, rawIV string) error {
	var err error

	iv, err = hex.DecodeString(rawIV)
	if err != nil {
		return err
	}

	tKey := strings.ReplaceAll(rawKey, "@", "")
	key, err = base64.URLEncoding.DecodeString(tKey)
	if err != nil {
		return err
	}

	for len(key) < 32 {
		key = append(key, 0)
	}

	block, err = aes.NewCipher(key)
	if err != nil {
		return err
	}

	return nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// EncryptByteArray encrypts a byte array with preset key/iv
func EncryptByteArray(data []byte) (string, error) {
	if block == nil {
		return "", errors.New("key/iv have not been set! Run InitCipher before attempting to encrypt/decrypt")
	}

	paddingData := pkcs5Padding(data, block.BlockSize())
	cipherData := make([]byte, len(paddingData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, paddingData)
	cipherDataAfter := base64.StdEncoding.EncodeToString(cipherData)
	return base64.URLEncoding.EncodeToString([]byte(cipherDataAfter + "::" + string(iv))), nil

}

// EncryptString wraps EncryptByteArray, expecting a string
func EncryptString(data string) (string, error) {
	return EncryptByteArray([]byte(data))
}

// DecryptString decrypts a string with preset key/iv
func DecryptString(data string) (string, error) {
	if block == nil {
		return "", errors.New("key/iv have not been set! Run InitCipher before attempting to encrypt/decrypt")
	}

	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	dataSplit := strings.Split(string(decoded), "::")
	if len(dataSplit) != 2 {
		return "", errors.New("could not split data")
	}

	decodedData := []byte(dataSplit[0])
	iv2 := []byte(dataSplit[1])

	mode := cipher.NewCBCDecrypter(block, iv2)

	decodedData2, err := base64.StdEncoding.DecodeString(string(decodedData))
	if err != nil {
		return "", err
	}

	for len(decodedData2)%block.BlockSize() != 0 {
		decodedData2 = append(decodedData2, 0)
	}

	mode.CryptBlocks(decodedData2, decodedData2)

	compiledString, err := regexp.Compile(`[^a-zA-Z0-9 -]`)
	if err != nil {
		return "", err
	}

	return compiledString.ReplaceAllString(string(decodedData2), ""), nil
}

// DecryptByteArray wraps Decrypt String, expecting a byte array
func DecryptByteArray(data []byte) (string, error) {
	return DecryptString(string(data))
}
