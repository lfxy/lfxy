package main

import (
	"fmt"
    "net/http"
	"io/ioutil"
	"strings"
	"time"
	"crypto/tls"
	"encoding/json"
	"bytes"
)

type UserDefinedMetricType string

const (
    // Instantaneous value. May increase or decrease.
    MetricGauge UserDefinedMetricType = "gauge"

    // A counter-like value that is only expected to increase.
    MetricCumulative UserDefinedMetricType = "cumulative"

    // Rate over a time period.
    MetricDelta UserDefinedMetricType = "delta"
)

type UserDefinedMetricDescriptor struct {
    // The name of the metric.
    Name string `json:"name"`

    // Type of the metric.
    Type UserDefinedMetricType `json:"type"`

    // Display Units for the stats.
    Units string `json:"units"`

    // Metadata labels associated with this metric.
    Labels map[string]string `json:"labels,omitempty"`
}

type UserDefinedMetric struct {
    UserDefinedMetricDescriptor `json:",inline"`
    // The time at which these stats were updated.
    Time time.Time `json:"time"`
    // Value of the metric. Float64s have 53 bit precision.
    // We do not foresee any metrics exceeding that value.
    Value float64 `json:"value"`
}

type UserDefinedMetricItems struct {
	UserDefinedMetrics map[string]UserDefinedMetric `json:"userDefinedMetrics,omitmepty"`
}

func httpGet(client http.Client,url string) {
    //resp, err := http.Get(url)
	resp, err := client.Get(url)
	//pc, file, line, ok := runtime.caller()
    if err != nil {
        // handle error
		return
    }

    body, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    fmt.Println(string(body))
    if err != nil {
        // handle error
    }
	//value := &UserDefinedMetricItems{ }
	/*var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	for k, v := range result {
		fmt.Println(k)
		fmt.Println(v)
		if v1, ok := v.(map[string]interface{}); ok {
			for k2, v2 := range v1 {
				fmt.Println("33333")
				fmt.Println(k2)
				fmt.Println(v2)
				if v3, ok := v2.([]interface{}); ok {
					for _, v4 := range v3 {
						fmt.Println("44444")
						fmt.Println(v4)
						if v5, ok := v4.(map[string]interface{}); ok {
							for k6, v6 := range v5 {
								fmt.Println("666666")
								fmt.Println(k6)
								fmt.Println(v6)

							}
						}
					}

				}
			}

		}
	}*/

	fmt.Println("_______________--------")
	var result map[string]map[string][]interface{}
	fmt.Errorf("failed to parse output. Response: %q.", string(body))
	err = json.Unmarshal(body, &result)
	for _, v := range result {
		for _, v1 := range v {
			for _, v2 := range v1 {
				if v3, ok := v2.(map[string]interface{}); ok {
					servername, s_exists := v3["server"];
					requestCounter, r_exists := v3["requestCounter"];
					if  s_exists && r_exists {
						fmt.Println(servername)
						if counter, ok := requestCounter.(float64); ok{
							fmt.Println(counter)
						}
					}
				}
			}

		}
	}
	fmt.Println("_______________--------")
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
	tNow := time.Now()
    fmt.Println(tNow)
	httpclient := createExternalHttpClient()
	httpGet(httpclient, "http://127.0.0.1:8989/qps")
	//httpGet(httpclient, "http://127.0.0.1:8989/qps/payss")
	//var body = "[{\"ip\":\"10.213.130.244:80\",\"cps\":\"50\"}, {\"ip\":\"10.213.130.245:80\",\"cps\":\"60\"}]"
	var body []map[string]string
	t1 := make(map[string]string)
	t1["ip"] = "10.213.130.244:80"
	t1["cps"] = "50"
	t2 := make(map[string]string)
	t2["ip"] = "10.213.130.245:80"
	t2["cps"] = "60"
	body = append(body, t1)
	body = append(body, t2)
	var result []map[string]string
	b, _ := json.Marshal(body)
	cnnn := string(b)
	fmt.Println(cnnn)
	fmt.Println(len(string(b)))
	err := json.Unmarshal(b, &result)
	if err != nil {
		fmt.Println("error;%s", err)
		return
	}
	for _, v := range result{
		for k1, v1 := range v {
			fmt.Println(k1)
			fmt.Println(v1)
		}
	}

	fmt.Println(bytes.Count(b, nil) - 1)
	str:="hhh"
	fmt.Println(strings.Count(str,"") - 1)
	fmt.Println(bytes.Count([]byte(str),nil)-1)
}

