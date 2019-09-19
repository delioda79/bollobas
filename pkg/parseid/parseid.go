package parseid

import (
	"bollobas/pkg/ciphrest"
	"regexp"
)

func main() {
	//fmt.Println("dfdf")
}

func EncryptString() {

}

func DecryptString(id string) string {

	//Here be temp cipher code..
	encodedID := regexp.MustCompile(`(?m).[^-]*.[^-]*$`).ReplaceAllString(id, "")

	//fmt.Printf("%s: %s --> %s\n", id, ciphrest.DecryptString(encodedID), ciphrest.EncryptString(ciphrest.DecryptString(encodedID)))

	return ciphrest.DecryptString(encodedID)

	//fmt.Println("decrypted", ciphrest.DecryptString(encodedID))
	//fmt.Println("re-encrypt", ciphrest.EncryptString(ciphrest.DecryptString(encodedID)))
}
