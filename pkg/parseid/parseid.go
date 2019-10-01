package parseid

import (
	"bollobas/pkg/ciphrest"
	"fmt"
	"os"
	"regexp"
)

// EncryptString encrypts an id with format: encryptedID-location-ut
func EncryptString(id int, ut string) string {
	location := os.Getenv("BOLLOBAS_LOCATION")
	if location == "" {
		panic("no location provided")
	}

	strID := fmt.Sprintf("%d", id)
	return fmt.Sprintf("%s-%s-%s", ciphrest.EncryptString(strID), location, ut)
}

// DecryptString decrypts id of a formated string of: encryptedID-location-ut
func DecryptString(id string) string {
	//Extracts encoded id (ignores location-pa/dr suffix)
	encodedID := regexp.MustCompile(`(?m).[^-]*.[^-]*$`).ReplaceAllString(id, "")

	return ciphrest.DecryptString(encodedID)
}
