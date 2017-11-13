package test

import (
	"context"
	"testing"

	"g.ghn.vn/go-common/dns-encapsulated/encapsulated"

	"github.com/stretchr/testify/assert"
)

var enscap *encapsulated.EncapsulatedConfig

type ClassDemo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func init() {
	enscap = encapsulated.DefaultConfig
	enscap.Debug = false
	enscap.BaseURI = "http://127.0.0.1:9898"

}
func TestGetHealth(t *testing.T) {
	resp, err := enscap.GetWithoutJson(context.Background(), "/health", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp, "true")

}
