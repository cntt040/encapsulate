package main

import (
	"net/http/httptest"

	"github.com/labstack/echo"

	"g.ghn.vn/go-common/dns-encapsulated/consul"
	"g.ghn.vn/go-common/encapsulate/encapsulated"
)

var log = encapsulated.GetLogger("Encapsulated with consul discover service")

func main() {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/health", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	consul := consul.Service
	consul.Debug = false
	consul.BaseUrlConsul = "http://192.168.101.213"
	consul.BasePortConsul = "53"

	type ClassDemo struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// GET
	resp1, err := consul.GetWithoutJson(c, "/health", nil)
	if err != nil {
		log.Error(err)
	}
	log.Info(resp1)

	//POST
	var resp2 ClassDemo
	err = consul.Post(c, "/health", nil, &resp2)
	if err != nil {
		log.Error(err)
	}
	log.Info(resp2)
}
