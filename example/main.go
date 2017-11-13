package main

import (
	"context"

	"g.ghn.vn/go-common/dns-encapsulated/encapsulated"
)

var log = encapsulated.GetLogger("Encapsulated Services")

func main() {
	enscap := encapsulated.DefaultConfig
	enscap.Debug = false
	enscap.BaseURI = "http://127.0.0.1:9898"

	type ClassDemo struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// GET
	resp1, err := enscap.GetWithoutJson(context.Background(), "/health", nil)
	if err != nil {
		log.Error(err)
	}
	log.Info(resp1)

	//POST
	var resp2 ClassDemo
	err = enscap.Post(context.Background(), "/health", nil, &resp2)
	if err != nil {
		log.Error(err)
	}
	log.Info(resp2)
}
