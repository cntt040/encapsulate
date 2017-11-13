package encapsulated

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type EncapsulatedConfig struct {
	BaseURI    string
	httpClient *http.Client
	Debug      bool
}

var DefaultConfig *EncapsulatedConfig

var logger = GetLoggerHTTP("Encapsulated Services")

func init() {
	def := &EncapsulatedConfig{
		BaseURI:    "",
		httpClient: http.DefaultClient,
		Debug:      true,
	}
	DefaultConfig = def
}
func (c *EncapsulatedConfig) GetWithoutJson(ctx echo.Context, path string, reqBody interface{}) (interface{}, error) {
	t0 := time.Now()
	reqData := c.encodeRequest(reqBody)
	req, err := http.NewRequest(echo.GET, c.BaseURI+path, bytes.NewBuffer(reqData))
	if err != nil {
		log.Panic(err)
		return nil, WrapError(err, CodeInternal, "Internal error")
	}

	//req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", ctx.Request().Header.Get("Authorization"))
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error")
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error, unable to read body")
	}

	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	if c.Debug == true {
		log.Infof("-> %s, st=%d, latency=%s, resp=%s", c.BaseURI+path, httpResp.StatusCode, t1.Sub(t0), string(respData))
	}

	return string(data), nil
}

func (c *EncapsulatedConfig) PostWithoutJson(ctx context.Context, path string, reqBody interface{}) (interface{}, error) {
	t0 := time.Now()
	reqData := c.encodeRequest(reqBody)
	req, err := http.NewRequest("POST", c.BaseURI+path, bytes.NewBuffer(reqData))
	if err != nil {
		log.Panic(err)
		return nil, WrapError(err, CodeInternal, "Internal error")
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error")
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error, unable to read body")
	}

	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	if c.Debug == true {
		log.Infof("-> %s, st=%d, latency=%s, resp=%s", c.BaseURI+path, httpResp.StatusCode, t1.Sub(t0), string(respData))
	}

	return string(data), nil
}

func (c *EncapsulatedConfig) Post(ctx context.Context, path string, reqBody interface{}, resp interface{}) error {
	t0 := time.Now()
	reqData := c.encodeRequest(reqBody)
	req, err := http.NewRequest("POST", c.BaseURI+path, bytes.NewBuffer(reqData))
	if err != nil {
		log.Panic(err)
		return WrapError(err, CodeInternal, "Internal error")
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return WrapError(err, CodeNetWork, "Network error")
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return WrapError(err, CodeNetWork, "Network error, unable to read body")
	}

	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	if c.Debug == true {
		log.Infof("-> %s, st=%d, latency=%s, resp=%s", c.BaseURI+path, httpResp.StatusCode, t1.Sub(t0), string(respData))
	}

	if resp == nil {
		resp = &Error{}
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return WrapError(err, CodeInternal, "Protocol unmarshal error "+string(respData))
	}

	return nil
}

func (c *EncapsulatedConfig) Get(ctx context.Context, path string, reqBody interface{}, resp interface{}) error {
	t0 := time.Now()
	reqData := c.encodeRequest(reqBody)
	req, err := http.NewRequest("GET", c.BaseURI+path, bytes.NewBuffer(reqData))
	if err != nil {
		log.Panic(err)
		return WrapError(err, CodeInternal, "Internal error")
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return WrapError(err, CodeNetWork, "Network error")
	}

	// Decode response
	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return WrapError(err, CodeNetWork, "Network error, unable to read body")
	}

	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	if c.Debug == true {
		log.Infof("-> %s, st=%d, latency=%s, resp=%s", c.BaseURI+path, httpResp.StatusCode, t1.Sub(t0), string(respData))
	}

	if resp == nil {
		resp = &Error{}
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return WrapError(err, CodeInternal, "Protocol unmarshal error "+string(respData))
	}

	return nil
}

func (c *EncapsulatedConfig) encodeRequest(reqBody interface{}) []byte {
	buf := new(bytes.Buffer)
	if reqBody == nil {
		return nil
	}

	err := json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		log.Panic(err)
	}

	b := buf.Bytes()
	if len(b) == 0 || b[0] != '{' {
		log.Panic(buf)
	}

	b = b[:len(b)-2]
	buf = bytes.NewBuffer(b)
	if len(b) > 1 {
		buf.WriteByte(',')
	}
	//buf.Write(c.cachedFields[1:])
	return buf.Bytes()
}
