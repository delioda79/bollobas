// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cipher_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestTralala(t *testing.T) {
	//e := cbcEncrypt()
	//fmt.Printf("Encrypted String : %s\n", e)
	d, pars := cbcDecrypt("NWZ3R2R4TUNQMEwxUStXdlZ4WUFsQT09Ojq+9N6igSEHCUtHEgUw4QyR", "s0th1s1s0uRpR1v@r3k3Y")
	//fmt.Printf("Decrypted String : %s\n", d)

	t.Logf("\nd: %s\n\n", string(d))
	t.Logf("\npars: %s\n\n", pars)
	t.Logf("\nencrypted_data: %s\n\n", pars[0])
	t.Logf("\nencryption_key: %s\n\n", pars[1])
	t.Logf("\ninit_vector: %s\n\n", pars[2])
	//fmt.Printf("%s\n", string(d))
	hi := string(d)
	_ = hi
	t.Error("yo")
}

func cbcDecrypt(data string, key string) ([]byte, [][]byte) {
	dataClean := strings.ReplaceAll(data, "@", "")
	dataDecoded, err := base64.StdEncoding.DecodeString(dataClean)
	if err != nil {
		panic(err)
	}

	keyClean := strings.ReplaceAll(key, "@", "")
	keyDecoded, err := base64.StdEncoding.DecodeString(keyClean)
	if err != nil {
		panic(err)
	}

	//APPEND
	keyDecoded = append(keyDecoded, 0)

	dataSplit := strings.Split(string(dataDecoded), "::")
	if len(dataSplit) != 2 {
		panic(errors.New("could not split data"))
	}

	dataRaw := []byte(dataSplit[0])
	initVector := []byte(dataSplit[1])

	//APPEND
	for i := 0; i < 8; i++ {
		dataRaw = append(dataRaw, 0)
	}

	//lala2 := fmt.Sprintf("encrypted_data: %s, encryption_key: %s, init_vecotr: %s\n", string(dataRaw), string(keyDecoded), string(initVector))

	return NewCBCDecrypter([]byte("1234567890123456"), []byte("1234567890123456"), []byte("1234567890123456")), [][]byte{dataRaw, keyDecoded, initVector}
	//return NewCBCDecrypter(dataRaw, keyDecoded, initVector), [][]byte{dataRaw, keyDecoded, initVector}
	//fmt.Printf("Decrypted String[] : %s\n", dStr)
	//nullC := "\000"

	//fmt.Printf("cKeyRaw : %s\n", cKeyRaw)

	//cKey2O, _ := base64.StdEncoding.DecodeString("WTd2Y0tmYTNyMWJqckpYQks1VzNLQT09OjqktsnVpjNYoy+Bj5CfWZee")
	//fmt.Printf("cKey2O : %s\n", cKey2O)

	/*
				block, err := aes.NewCipher(key)
				if err != nil {	lala2 := fmt.Sprintf("encrypted_data: %s, encryption_key: %s, init_vecotr: %s\n", dataRaw, key	lala2 := fmt.Sprintf("encrypted_data: %s, encryption_key: %s, init_vecotr: %s\n", dataRaw, keyDecoded, initVector)
			_ = lala2
			return NewCBCDecrypter(dataRaw, keyDecoded, initVector)
			//fmt.Printf("Decrypted String[] : %s\n", dStr)
			lala2 := fmt.Sprintf("encrypted_data: %s, encryption_key: %s, init_vecotr: %s\n", dataRaw, keyDecoded, initVector)
			_ = lala2
			return NewCBCDecrypter(dataRaw, keyDecoded, initVector)
			//fmt.Printf("Decrypted String[] : %s\n", dStr)
		Decoded, initVector)
			_ = lala2
			return NewCBCDecrypter(dataRaw, keyDecoded, initVector)
			//fmt.Printf("Decrypted String[] : %s\n", dStr)

					panic(err)
				}

				// include it at the beginning of the ciphertext.
				if len(ciphertext) < aes.BlockSize {
					panic("ciphertext too short")
				}
				iv := ciphertext[:aes.BlockSize]
				ciphertext = ciphertext[aes.BlockSize:]

				// CBC mode always works in whole blocks.
				if len(ciphertext)%aes.BlockSize != 0 {
					panic("ciphertext is not a multiple of the block size")
				}

				mode := cipher.NewCBCDecrypter(block, iv)

				// CryptBlocks can work in-place if the two arguments are the same.
				mode.CryptBlocks(ciphertext, ciphertext)
				ciphertext = PKCS5UnPadding(ciphertext)
				return ciphertext*/
}

func NewCBCDecrypter(ciphertext []byte, key []byte, iv []byte) []byte {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	/*
		key, _ := hex.DecodeString("6368616e676520746869732070617373")
		ciphertext, _ := hex.DecodeString("73c86d43a9d700a253a96c85b0f6b03ac9792e0e757f869cca306bd3cba1c62b")
	*/
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	/*
		// The IV needs to be unique, but not secure. Therefore it's common to
		// include it at the beginning of the ciphertext.
		if len(ciphertext) < aes.BlockSize {
			panic("ciphertext too short")
		}
		iv := ciphertext[:aes.BlockSize]
		ciphertext = ciphertext[aes.BlockSize:]
	*/
	lala := aes.BlockSize
	_ = lala
	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	return ciphertext
}

func ExampleNewCBCEncrypter() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Printf("%x\n", ciphertext)
}

func cbcEncrypt() string {
	key := []byte("keyforencryption")
	plaintext := []byte("testssssss")
	plaintext = PKCS5Padding(plaintext, 16)
	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
