package parseid

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/taxibeat/bollobas/pkg/ciphrest"
)

// EncryptString encrypts an id with format: encryptedID-location-ut
func EncryptString(id int, ut string) (string, error) {
	location := os.Getenv("BOLLOBAS_LOCATION")
	if location == "" {
		return "", errors.New("no location provided")
	}

	strID := fmt.Sprintf("%d", id)
	encrypted, err := ciphrest.EncryptString(strID)
	if err != nil {

	}
	return fmt.Sprintf("%s-%s-%s", encrypted, location, ut), nil
}

// DecryptString decrypts id of a formated string of: encryptedID-location-ut
func DecryptString(id string) (string, error) {
	//Extracts encoded id (ignores location-pa/dr suffix)
	compiledString, err := regexp.Compile(`(?m).[^-]*.[^-]*$`)
	if err != nil {
		return "", err
	}

	encodedID := compiledString.ReplaceAllString(id, "")
	return ciphrest.DecryptString(encodedID)
}
