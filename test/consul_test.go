package test

import (
	"net/http/httptest"
	"testing"

	"g.ghn.vn/go-common/dns-encapsulated/consul"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthConsul(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/health", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	var consul consul.ClientDns
	consul.Debug = true
	consul.BaseUrlConsul = "192.168.101.213"
	consul.BasePortConsul = "53"
	consul.NameService = "station_hydros-9898"
	consul.InitService()

	resp, err := consul.GetWithoutJson(c, "/health", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp, "true")

}
