package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfiguration(t *testing.T) {
	cfg := &Configuration{}
	_, e := NewConfig(cfg)

	assert.Nil(t, e)
}
