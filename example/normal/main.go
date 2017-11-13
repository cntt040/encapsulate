package main

import (
	"encoding/json"
	"net/http/httptest"

	"github.com/labstack/echo"

	"g.ghn.vn/go-common/dns-encapsulated/encapsulated"
)

var log = encapsulated.GetLogger("Encapsulated Services")

func main() {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/health", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	enscap := encapsulated.DefaultConfig
	enscap.Debug = false
	enscap.BaseURI = "http://127.0.0.1:9898"

	type ClassDemo struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// GET
	resp1, err := enscap.RequestWithoutJson(c, echo.GET, "/health", nil)
	if err != nil {
		log.Error(err)
	}
	log.Info(resp1)

	//POST
	var resp2 ClassDemo
	data, err := enscap.Request(c, echo.POST, "/health", nil)
	if err != nil {
		log.Error(err)
	}
	err = json.Unmarshal(data, &resp2)
	if err != nil {
		log.Error(err)
	}

	log.Info(resp2)
}
