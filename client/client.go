package main

import (
	"fmt"
    "net/http"
	"io/ioutil"
	//"strings"
	//"time"
	"crypto/tls"
	"encoding/json"
	"runtime/debug"
)


func httpGet(client http.Client,url string) {
    //resp, err := http.Get(url)
	resp, err := client.Get(url)
	//pc, file, line, ok := runtime.caller()
	debug.PrintStack()
    if err != nil {
        // handle error
		return
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // handle error
    }

    fmt.Println(string(body))
}

func createExternalHttpClient() http.Client {
   tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
    }

    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }

    return http.Client{Transport: transport}
}





func main() {
	/*metricName := "metrics/qps"
	tNow := time.Now()
    fmt.Println(tNow)
//	if strings.Contains(metricName, "qps") {
	if metricName == "metrics/qps" {
		fmt.Println("aaaaaa")
	}
	httpclient := createExternalHttpClient()
	httpGet(httpclient, "http://127.0.0.1:8989/qps/payss")
	//httpGet(httpclient, "http://127.0.0.1:8989/qps/payss")
    */
    p_str := "{\"c08b5377-4cdb-4134-90d6-9084a37f538b\":1}"
    p_m := make(map[string]int)
    err := json.Unmarshal([]byte(p_str), &p_m)
    if err != nil {
        fmt.Printf("err:%s", err.Error())
    }
    fmt.Println(p_m)
}

