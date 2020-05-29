package main

import (
    "fmt"
	"net/http"
	"io/ioutil"
	"strings"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/push"
)


func httpPost(text string) {
    resp, err := http.Post("http://heapster/api/v1/push/prometheus",
        "application/x-www-form-urlencoded",
        strings.NewReader(text))
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		fmt.Println(err)
        // handle error
    }
    fmt.Println(string(body))
}

func main() {
    completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "http_requests_per_minute",
        Help: "The timestamp of the last successful completion of a DB backup.",
		ConstLabels: prometheus.Labels{"namespace": "kube-system", "pod": "php-apache-9h7hk"},
    })

	/*text := `
# This is a pod-level metric (it might be used for autoscaling)
# TYPE http_requests_per_second guage
http_requests_per_minute{namespace="kube-system",pod="php-apache-9h7hk"} 20
http_requests_per_minute{namespace="kube-system",pod="php-apache-g7s4q"} 5

# This is a service-level metric, which will be stored as frontend_hits_total
# and restapi_hits_total (these might be used for auto-idling)
# TYPE hits_total counter
hits_total{namespace="kube-system",service="frontend"} 5000
hits_total{namespace="kube-system",service="nginxservice"} 6000
`*/
	text2 := `http_requests_per_minute{namespace="kube-system",pod="php-apache-9h7hk"} 20`
	httpPost(text2)
	completionTime.Set(1000)
    //completionTime.SetToCurrentTime()
	//reg := prometheus.NewRegistry()
	//reg.MustRegister(completionTime)
    if err := push.AddCollectors(
        "http_gatherer", push.HostnameGroupingKey(),
        "http://heapster/api/v1/push/prometheus/",
        completionTime,
    ); err != nil {
        fmt.Println("Could not push completion time to Pushgateway:", err)
    }
}
