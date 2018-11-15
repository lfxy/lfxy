package main

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
)

func test_get() error {
	request := gorequest.New()
    resp, body, errs := request.Get("http://localhost:8081/api/v1/cluster/all").End()
    fmt.Printf("resp:%s\n", resp)
    fmt.Printf("body:%s\n", body)
    fmt.Printf("err:%s\n", errs)
    return nil
}

func main(){
    test_get()
}
