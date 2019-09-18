package parseid

import (
	"bollobas/pkg/ciphrest"
	"fmt"
	"regexp"
)

func main() {
	fmt.Println("dfdf")
}

func EncryptString() {

}

func DecryptString(id string) string {
	//Here be temp cipher code..
	encodedID := regexp.MustCompile(`(?m).[^-]*.[^-]*$`).ReplaceAllString(id, "")
	return ciphrest.DecryptString(encodedID)
	//fmt.Println("original", encodedID)
	//fmt.Println("decrypted", ciphrest.DecryptString(encodedID))
	//fmt.Println("re-encrypt", ciphrest.EncryptString(ciphrest.DecryptString(encodedID)))
}
