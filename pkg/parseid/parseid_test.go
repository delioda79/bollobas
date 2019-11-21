package parseid

import (
	"fmt"
	"os"
	"testing"

	"github.com/taxibeat/bollobas/pkg/ciphrest"

	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	err := ciphrest.InitCipher("44441s111111R1222221", "11111111112222222222333333333344")
	assert.NoError(t, err)
	type expectation struct {
		errors assert.ErrorAssertionFunc
	}
	type args struct {
		input string
	}
	tests := map[string]struct {
		args        args
		expectation expectation
	}{
		"compileable-encrypted-ID": {
			args: args{
				input: "MjZtQmRsZEgrdEQyOHJrYkI4UkczQT09OjoRERERESIiIiIiMzMzMzNE-sandbox-pa",
			},
			expectation: expectation{
				errors: assert.NoError,
			},
		},
		"non-compileable-ID": {
			args: args{
				input: "1",
			},
			expectation: expectation{
				errors: assert.Error,
			},
		},
		"compileable-non-encrypted-id": {
			// follows proper regex but can't be decrypted
			args: args{
				input: "1-sandbox-pa",
			},
			expectation: expectation{
				errors: assert.Error,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dc, err := DecryptString(tt.args.input)
			tt.expectation.errors(t, err)
			if err == nil {
				assert.Equal(t, fmt.Sprintf("%d", 4104), dc)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	os.Setenv("BOLLOBAS_LOCATION", "sandbox")
	err := ciphrest.InitCipher("44441s111111R1222221", "11111111112222222222333333333344")
	assert.NoError(t, err)
	dc, err := EncryptString(4104, "pa")
	assert.NoError(t, err)
	assert.Equal(t, "MjZtQmRsZEgrdEQyOHJrYkI4UkczQT09OjoRERERESIiIiIiMzMzMzNE-sandbox-pa", dc)
}
