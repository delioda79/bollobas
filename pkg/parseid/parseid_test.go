package parseid

import (
	"fmt"
	"os"
	"testing"

	"github.com/taxibeat/bollobas/pkg/ciphrest"

	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	ciphrest.InitCipher("44441s111111R1222221", "11111111112222222222333333333344")
	dc := DecryptString("MjZtQmRsZEgrdEQyOHJrYkI4UkczQT09OjoRERERESIiIiIiMzMzMzNE-sandbox-pa")
	assert.Equal(t, fmt.Sprintf("%d", 4104), dc)
}

func TestEncrypt(t *testing.T) {
	os.Setenv("BOLLOBAS_LOCATION", "sandbox")
	ciphrest.InitCipher("44441s111111R1222221", "11111111112222222222333333333344")
	dc := EncryptString(4104, "pa")
	assert.Equal(t, "MjZtQmRsZEgrdEQyOHJrYkI4UkczQT09OjoRERERESIiIiIiMzMzMzNE-sandbox-pa", dc)
}
