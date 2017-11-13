package test

import (
	"net/http/httptest"
	"testing"

	"g.ghn.vn/go-common/dns-encapsulated/encapsulated"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var enscap *encapsulated.EncapsulatedConfig

type ClassDemo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func init() {
	enscap = encapsulated.DefaultConfig
	enscap.Debug = true
	enscap.BaseURI = "http://10.100.144.150:9898"

}
func TestGetHealth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/health", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	resp, err := enscap.GetWithoutJson(c, "/health", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp, "true")

}

