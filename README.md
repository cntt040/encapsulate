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
	err = enscap.Post(context.Background(), "/health", nil, &resp2)
	if err != nil {
		log.Error(err)
	}
    ```

    Func get api as REST APIs