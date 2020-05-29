package main
import (
	"bytes"
	"encoding/json"
    "fmt"
    "net/http"
    "strings"
    "log"
	"time"
	io "io/ioutil"
	//"os"
)

var count = 0

type CustomMetricTarget struct {
    // Custom Metric name.
    Name string `json:"name"`
    // Custom Metric value (average).
    TargetValue int `json:"value"`
    ServiceName string `json:"servicename,omitempty"`
}

type CustomMetricTargetList struct {
	Items []CustomMetricTarget `json:"items"`

}

type CustomResult struct {
	Timestamp time.Time
	Value int64
}


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
    //Labels map[string]string `json:"labels,omitempty"`
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
	UserDefinedMetrics map[string][]UserDefinedMetric `json:"userDefinedMetrics,omitmepty"`
}

func (uidMetaDataResp *UserDefinedMetricItems) String() string {
    buffer := bytes.NewBuffer(nil)
    content, _ := json.Marshal(uidMetaDataResp)
    buffer.WriteString(fmt.Sprintf("%s\n", string(content)))
    return buffer.String()
}


func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  //解析参数，默认是不会解析的
    //fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	//jstr := `{"timestamp":"2016-12-23T12:08:00Z","value":94}`
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
	now := time.Now()
	udmd := UserDefinedMetricDescriptor {
		"qps",
		MetricGauge,
		"number of requests",
	}
	udm := []UserDefinedMetric {
		{udmd, now, 90, },
	}

	udms := make(map[string][]UserDefinedMetric)
	udms["service1"] = udm

	um1 := &UserDefinedMetricItems {
		udms,
	}

	b2, err := json.Marshal(udms)
	if err != nil {
		fmt.Println("error:", err)
    } else {
        fmt.Println("encoded data : ")
        fmt.Println(string(b2))
    }
	fmt.Println(um1.String())
	fmt.Println("------------------")


	data, _ := io.ReadFile("./data1.json")

	//fmt.Println(data)

	fmt.Println("------------------")

	var result map[string]map[string][]interface{}
	err = json.Unmarshal(data, &result)
	for _, v := range result {
		for _, v1 := range v {
			for _, v2 := range v1 {
				if v3, ok := v2.(map[string]interface{}); ok {
					servername, s_exists := v3["server"];
					_, r_exists := v3["requestCounter"];
					if  s_exists && r_exists {
						fmt.Println(servername)
						if counter, ok2 := v3["requestCounter"].(float64); ok2{
							v3["requestCounter"] = counter + float64(55 * 60 * count)
							fmt.Println("czq:%f, %d", v3["requestCounter"], count)
						}
					}
				}
			}
		}
	}
	count++
	data2, err := json.Marshal(result)
	//os.Stdout.Write(b)
    fmt.Fprintf(w, string(data2)) //这个写入到w的是输出到客户端的
}

func main() {

    http.HandleFunc("/qps", sayhelloName) //设置访问的路由
    err := http.ListenAndServe(":8990", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
