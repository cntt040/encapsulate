** ENCAPSULATED DNS 
    ```
    git clone git@g.ghn.vn:go-common/encapsulate.git
    ```

** Test
    ```
    go test ./test/...
    ```

*** Update
    ```
    go get -u all 
    go get -u g.ghn.vn/go-common/encapsulate
    ```
*** Example

Create servive 

    ```
    enscap := encapsulated.DefaultConfig
	enscap.Debug = false
	enscap.BaseURI = "http://127.0.0.1:9898"
    ```

    enscap.BaseURI = "http://127.0.0.1:9898" : base url

    ```
    resp1, err := enscap.GetWithoutJson(context.Background(), "/health", nil)
	if err != nil {
		log.Error(err)
	}
    ```
    If you need get data from http as string, you can use func GetWithoutJson or PostWithoutJson

    ```
    var resp2 ClassDemo
	data, err := enscap.Request(c, echo.POST, "/health", nil)
	if err != nil {
		log.Error(err)
	}
	err = json.Unmarshal(data, &resp2)
	if err != nil {
		log.Error(err)
	}
    ```

Create service with consul
    ```
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
    ```

    ```
    === RUN   TestGetHealthConsul
    2017/11/13 22:43:13 [INFO] &{station_hydros-9898 192.168.101.213 53 .service.consul true <nil>}
    *22:43:13.820 -> http://10.100.144.150:9896/health, st=200, latency=14.833453ms, resp=true
    --- PASS: TestGetHealthConsul (0.03s)
    === RUN   TestGetHealth
    *22:43:13.835 -> http://10.100.144.150:9898/health, st=200, latency=15.215197ms, resp=true
    --- PASS: TestGetHealth (0.02s)
    PASS
    ok      g.ghn.vn/go-common/dns-encapsulated/test        0.048s

    ```
    Func get api as REST APIs