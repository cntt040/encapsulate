package encapsulated

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type EncapsulatedConfig struct {
	BaseURI    string
	httpClient *http.Client
	Debug      bool
}

var DefaultConfig *EncapsulatedConfig

var logger = GetLoggerHTTP("Encapsulated Services ")

func init() {
	def := &EncapsulatedConfig{
		BaseURI:    "",
		httpClient: http.DefaultClient,
		Debug:      true,
	}
	DefaultConfig = def
}

func (c *EncapsulatedConfig) RequestWithoutJson(ctx echo.Context, method string, path string, reqBody interface{}) (interface{}, error) {
	t0 := time.Now()
	reqData := c.encodeRequest(reqBody)
	req, err := http.NewRequest(method, c.BaseURI+path, bytes.NewBuffer(reqData))
	if err != nil {
		logger.Panic(err)
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
		logger.Infof("-> %s, st=%d, latency=%s,req=%s, resp=%s", c.BaseURI+path, httpResp.StatusCode, t1.Sub(t0), string(reqData), string(respData))
	}
	if httpResp.StatusCode >= 300 || httpResp.StatusCode < 200 {
		var resErr *Error
		es := json.Unmarshal(respData, &resErr)
		if es != nil || resErr == nil {
			return nil, WrapError(nil, strconv.Itoa(httpResp.StatusCode), "Unmarshal response")
		}
		return nil, WrapError(nil, resErr.Code, resErr.Message)
	}
	return string(data), nil
}

func (c *EncapsulatedConfig) Request(ctx echo.Context, method string, path string, reqBody interface{}) ([]byte, error) {
	t0 := time.Now()
	reqData := c.encodeRequest(reqBody)
	req, err := http.NewRequest(method, c.BaseURI+path, bytes.NewBuffer(reqData))
	if err != nil {
		logger.Panic(err)
		return nil, WrapError(err, CodeInternal, "Internal error")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization", ctx.Request().Header.Get("Authorization"))
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error")
	}

	// Decode response
	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error, unable to read body")
	}

	if httpResp.StatusCode >= 300 || httpResp.StatusCode < 200 {
		var resErr *Error
		es := json.Unmarshal(data, &resErr)
		if es != nil || resErr == nil {
			fmt.Println("return 1")
			return nil, WrapError(es, strconv.Itoa(httpResp.StatusCode), "Unmarshal response")
		}
		fmt.Println("return 2")
		return nil, WrapErrorf(resErr.Code, resErr.Message)
	}

	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	if c.Debug == true {
		logger.Infof("-> %s, st=%d, latency=%s, resp=%s", c.BaseURI+path, httpResp.StatusCode, t1.Sub(t0), string(respData))
	}

	return data, nil
}

func (c *EncapsulatedConfig) encodeRequest(reqBody interface{}) []byte {
	buf := new(bytes.Buffer)
	if reqBody == nil {
		return nil
	}

	err := json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		logger.Panic(err)
	}

	b := buf.Bytes()
	if len(b) == 0 || b[0] != '{' {
		logger.Panic(buf)
	}

	b = b[:len(b)-1]
	buf = bytes.NewBuffer(b)
	// if len(b) > 1 {
	// 	buf.WriteByte(',')
	// }
	//buf.Write(c.cachedFields[1:])
	return buf.Bytes()
}
