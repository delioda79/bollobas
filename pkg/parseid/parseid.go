package parseid

import (
	"bollobas/pkg/ciphrest"
	"fmt"
	"os"
	"regexp"
)

func main() {
	//fmt.Println("dfdf")
}

func EncryptString(id int, ut string) string {
	location := os.Getenv("BOLLOBAS_LOCATION")
	if location == "" {
		panic("no location provided")
	}

	strID := fmt.Sprintf("%d", id)
	return fmt.Sprintf("%s-%s-%s", ciphrest.EncryptString(strID), location, ut)
}

func DecryptString(id string) string {

	//Here be temp cipher code..
	encodedID := regexp.MustCompile(`(?m).[^-]*.[^-]*$`).ReplaceAllString(id, "")

	//fmt.Printf("%s: %s --> %s\n", id, ciphrest.DecryptString(encodedID), ciphrest.EncryptString(ciphrest.DecryptString(encodedID)))

	return ciphrest.DecryptString(encodedID)

	//fmt.Println("decrypted", ciphrest.DecryptString(encodedID))
	//fmt.Println("re-encrypt", ciphrest.EncryptString(ciphrest.DecryptString(encodedID)))
}
