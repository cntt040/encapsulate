package consul

import (
	"fmt"

	"g.ghn.vn/go-common/dns-encapsulated/encapsulated"
	"github.com/benschw/dns-clb-go/clb"
	"github.com/labstack/echo"
	"github.com/yanzay/log"
)

type ClientDns struct {
	NameService    string
	BaseUrlConsul  string
	BasePortConsul string
	BasePathConsul string
	Debug          bool
	encap          *encapsulated.EncapsulatedConfig
}

func (c *ClientDns) InitService() {
	c.BasePathConsul = ".service.consul"
	enscap := encapsulated.DefaultConfig
	enscap.Debug = c.Debug

	uri, err := c.getAddress()
	if err != nil {
		uri = ""
	}
	log.Info(c)
	enscap.BaseURI = fmt.Sprintf("http://%s", uri)
	c.encap = enscap

}

func (c *ClientDns) getAddress() (string, error) {
	cs := clb.NewClb(c.BaseUrlConsul, c.BasePortConsul, clb.Random)
	srvRecord := c.NameService + c.BasePathConsul
	address, err := cs.GetAddress(srvRecord)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func (c *ClientDns) Get(ctx echo.Context, path string, reqBody interface{}) ([]byte, error) {
	res, err := c.encap.Get(ctx, path, reqBody)
	return res, err
}

func (c *ClientDns) Post(ctx echo.Context, path string, reqBody interface{}, resp interface{}) error {
	return c.encap.Post(ctx, path, reqBody, reqBody)
}

func (c *ClientDns) GetWithoutJson(ctx echo.Context, path string, reqBody interface{}) (interface{}, error) {
	res, err := c.encap.GetWithoutJson(ctx, path, reqBody)
	return res, err
}

func (c *ClientDns) PostWithoutJson(ctx echo.Context, path string, reqBody interface{}) (interface{}, error) {
	res, err := c.encap.PostWithoutJson(ctx, path, reqBody)
	return res, err
}
