package encapsulated

import (
	"bytes"
	"encoding/json"
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
		log.Info(DataLog{
			Url:    c.BaseURI + path,
			Status: httpResp.StatusCode,
			St:     t1.Sub(t0),
			Req:    string(reqData),
			Res:    string(respData),
		})
	}
	if httpResp.StatusCode >= 300 || httpResp.StatusCode < 200 {
		var resErr *Error
		es := json.Unmarshal(data, &resErr)
		if es != nil || resErr == nil {

			return nil, WrapError(es, strconv.Itoa(httpResp.StatusCode), "Unmarshal response")
		}

		return nil, WrapErrorf(resErr.Code, resErr.Message)
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

	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	if c.Debug == true {
		log.Info(DataLog{
			Url:    c.BaseURI + path,
			Status: httpResp.StatusCode,
			St:     t1.Sub(t0),
			Req:    string(reqData),
			Res:    string(respData),
		})
	}
	if httpResp.StatusCode >= 300 || httpResp.StatusCode < 200 {
		var resErr *Error
		es := json.Unmarshal(data, &resErr)
		if es != nil || resErr == nil {

			return nil, WrapError(es, strconv.Itoa(httpResp.StatusCode), "Unmarshal response")
		}

		return nil, WrapErrorf(resErr.Code, resErr.Message)
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

func (c *EncapsulatedConfig) RequestWithHeaer(ctx echo.Context, method string, path string, reqBody interface{}, header []string) ([]byte, error) {
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
	for _, h := range header {
		req.Header.Add(h, ctx.Request().Header.Get(h))
	}

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapError(err, CodeNetWork, "Network error")
	}

	// Decode response
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
		log.Info(DataLog{
			Url:    c.BaseURI + path,
			Status: httpResp.StatusCode,
			St:     t1.Sub(t0),
			Req:    "",
			Res:    string(respData),
		})
	}
	if httpResp.StatusCode >= 300 || httpResp.StatusCode < 200 {
		var resErr *Error
		es := json.Unmarshal(data, &resErr)
		if es != nil || resErr == nil {

			return nil, WrapError(es, strconv.Itoa(httpResp.StatusCode), "Unmarshal response")
		}

		return nil, WrapErrorf(resErr.Code, resErr.Message)
	}
	return data, nil
}
