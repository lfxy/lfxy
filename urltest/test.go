package main

import (
	"flag"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"
	"net/http"
	"crypto/tls"
	"io/ioutil"
)

type Uri struct {
	Key string
	Val url.URL
}

func (u *Uri) String() string {
	val := u.Val.String()
	if val == "" {
		return fmt.Sprintf("%s", u.Key)
	}
	return fmt.Sprintf("%s:%s", u.Key, val)
}

func (u *Uri) Set(value string) error {
	s := strings.SplitN(value, ":", 2)
	if s[0] == "" {
		return fmt.Errorf("missing uri key in '%s'", value)
	}
	u.Key = s[0]
	if len(s) > 1 && s[1] != "" {
		e := os.ExpandEnv(s[1])
		uri, err := url.Parse(e)
		if err != nil {
			return err
		}
		u.Val = *uri
	}
	return nil
}

type Uris []Uri

func (us *Uris) String() string {
	var b bytes.Buffer
	b.WriteString("[")
	for i, u := range *us {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(u.String())
	}
	b.WriteString("]")
	return b.String()
}

func (us *Uris) Set(value string) error {
	var u Uri
	if err := u.Set(value); err != nil {
		return err
	}
	*us = append(*us, u)
	return nil
}

var argCustoms Uris


func createExternalHttpClient() http.Client {
   tlsConfig := &tls.Config{
        InsecureSkipVerify: true,
    }

    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
    }

    return http.Client{Transport: transport}
}


func httpGet(client http.Client,url string) {
    //resp, err := http.Get(url)
	resp, err := client.Get(url)
	//pc, file, line, ok := runtime.caller()
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


func postRequestAndGetValue(client *http.Client, req *http.Request, rq string) error {
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
		fmt.Println("body: error")
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed - , response: ")
	}
	fmt.Println("body: ", string(body))
	return nil
}

func GetCustomMetrics(client *http.Client, host string, path string, rq string) (error) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   path,
		RawQuery:   rq,
	}
	fmt.Println("czq kubelet_client.go, GetCustomMetrics url:  ", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	if client == nil {
		client = http.DefaultClient
	}
	err = postRequestAndGetValue(client, req, rq)
	return err
}
func main() {
	flag.Var(&argCustoms, "custom", "custom metrics collect from")
	flag.Parse()
	fmt.Println(len(argCustoms))
	fmt.Println("string:", argCustoms[0].Val.String())
	fmt.Println("Scheme:", argCustoms[0].Val.Scheme)
	fmt.Println("Opaque:", argCustoms[0].Val.Opaque)
	fmt.Println("Path:", argCustoms[0].Val.Path)
	fmt.Println("RawPath:", argCustoms[0].Val.RawPath)
	fmt.Println("ForceQuery:", argCustoms[0].Val.ForceQuery)
	fmt.Println("RawQuery:", argCustoms[0].Val.RawQuery)
//	fmt.Println("Fragment:", argCustoms[0].Val.Fragment)
	httpclient := createExternalHttpClient()
	path := argCustoms[0].Val.Path// + "?" + argCustoms[0].Val.RawQuery
	GetCustomMetrics(&httpclient, argCustoms[0].Val.Host, path, argCustoms[0].Val.RawQuery)
}
