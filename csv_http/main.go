package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"crypto/tls"
	"time"
)

func createHaproxyClient() *http.Client {
   tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
    }

    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }

	c := &http.Client{
		Transport: transport,
		//Timeout:   kubeletConfig.HTTPTimeout,
	}
	return c
}

func GetReloadTime(client *http.Client, ip string, port int, path string) (map[string]time.Time, error) {
	//result := make(map[string]time.Time)
	result := map[string]time.Time{}
	url := url.URL {
		Scheme:		"http",
		Host:		fmt.Sprintf("%s:%d", ip, port),
		Path:		path,
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return result, err
	}
	if client == nil {
		client = http.DefaultClient
	}

	err = postReloadTimeReq(client, req, &result)

	return result, err
}

func postReloadTimeReq(client *http.Client, req *http.Request, result *map[string]time.Time) (error) {
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body - %v", err)
	}
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
	}

	//to do
	timearr := strings.Split(string(body), "\n")
	if len(timearr) != 2 {
		return fmt.Errorf("response body does not have right time")
	}

	lasttime, err := time.Parse("2006-01-02 15:04:05", timearr[0])
	if err != nil {
		return fmt.Errorf("parse last time error")
	}
	currenttime, err := time.Parse("2006-01-02 15:04:05", timearr[1])
	if err != nil {
		return fmt.Errorf("parse current time error")
	}
	fmt.Println(lasttime)
	fmt.Println("-------")
	fmt.Println(currenttime)
	(*result)["LastReloadTime"] = lasttime
	(*result)["CurrentTime"] = currenttime

	return nil
}
func main(){
	/*c := createHaproxyClient()
	cm, err := GetReloadTime(c, "localhost", 33966, "reloadtime")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cm)*/
	t1 := "2017-03-26 23:14:52.003711863 -0400 EDT"
	currenttime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", t1)
	if err != nil {
		fmt.Println(err)
		return
	}
	c2 := currenttime.Format("2006-01-02 15:04:05")
	fmt.Println(currenttime)
	fmt.Println(currenttime.String())
	fmt.Println(c2)

	m := make(map[string]map[int]float64)
	m["aaa"] = make(map[int]float64)
	m["aaa"][1] = 6.000000

	for k1, v1 := range m {
		fmt.Println(k1)
		for k2, v2 := range v1 {
			fmt.Println(k2, v2)
		}
	}
}
