package ciph

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var key []byte
var iv []byte
var block cipher.Block

func Init() {
	var err error

	iv, err = hex.DecodeString("bef4dea2812107094b47120530e10c91")
	if err != nil {
		panic(err)
	}

	//	fmt.Println("iv Length:", len(iv))
	tKey := strings.ReplaceAll("s0th1s1s0uRpR1v@r3k3Y", "@", "")
	key, err = base64.URLEncoding.DecodeString(tKey)
	if err != nil {
		panic(err)
	}

	for len(key) < 32 {
		key = append(key, 0)
	}

	block, err = aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Init")
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func main() {
	//fmt.Printf("BlockSize: %d\n", block.BlockSize())

	//	encrypted := EncryptToken(key, iv, []byte("fubar"))
	//	fmt.Printf("Go encrypt: %s\n", encrypted)
	//	fmt.Printf("Go decrypt: %s\n", DecryptToken(key, iv, encrypted))
	//var re = regexp.MustCompile(`(?m).[^-]*.[^-]*$`)
	//var str = `bXdFZjZVOTdIbFJtc1I1Ylo4QXNmZz09Ojq-9N6igSEHCUtHEgUw4QyR-sandbox-dr`

	//fmt.Println(regexp.MustCompile(`(?m).[^-]*.[^-]*$`).ReplaceAllString(str, ""))

	mData := "7040"
	fmt.Println("Data:", mData)
	encrypted := EncryptData([]byte(mData), key, iv)
	fmt.Printf("Go encrypt: %s\n", encrypted)
	//fmt.Println("Encrypt EXPECTED: b0pBMFpSVk1rNVI3anY1bWwwMURBQT09Ojq-9N6igSEHCUtHEgUw4QyR")
	//fmt.Printf("Go decrypt2: %s\n", DecryptData(encrypted2, key))
	fmt.Printf("Go decrypt: %s\n", DecryptData(encrypted))
}

func EncryptData(data []byte, key []byte, iv []byte) string {
	paddingData := PKCS5Padding(data, block.BlockSize())
	cipherData := make([]byte, len(paddingData))

	//fmt.Println("iv2", string(iv))
	//fmt.Println("lenCipher", len(cipherData))
	//fmt.Println("key", string(key))
	//fmt.Println("cipherData:", string(cipherData))
	//fmt.Println("paddingData:", string(paddingData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, paddingData)
	//fmt.Println("cipherDataAfter:", string(cipherData))
	cipherDataAfter := string(base64.URLEncoding.EncodeToString(cipherData))
	//fmt.Println("cipherDataAfterBase:", cipherDataAfter)
	return base64.URLEncoding.EncodeToString([]byte(cipherDataAfter + "::" + string(iv)))

}

func DecryptData(data string) string {
	//fmt.Println(data)
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Decodedbase64:", decoded, len(decoded))
	dataSplit := strings.Split(string(decoded), "::")
	if len(dataSplit) != 2 {
		panic(errors.New("could not split data"))
	}

	decodedData := []byte(dataSplit[0])
	iv2 := []byte(dataSplit[1])

	//fmt.Println("iv2", string(iv2))
	mode := cipher.NewCBCDecrypter(block, iv2)

	//fmt.Println("len", len(decodedData))
	//fmt.Println("key", string(key))
	//fmt.Println("decodedData:", string(decodedData))

	decodedData2, err := base64.URLEncoding.DecodeString(string(decodedData))
	if err != nil {
		panic(err)
	}
	//fmt.Println("DecodedData2:", string(decodedData2), len(decodedData2))
	for len(decodedData2)%block.BlockSize() != 0 {
		decodedData2 = append(decodedData2, 0)
	}

	mode.CryptBlocks(decodedData2, decodedData2)
	if err != nil {
		panic(err)
	}
	return string(decodedData2)
}
