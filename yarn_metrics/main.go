package main

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
    "os"
    "encoding/json"
	//"sort"
    "encoding/csv"
    "flag"
    "time"
    //"strings"
    "strconv"
)

type RMMetrics struct {
    MetricsObj         MetricsInfo          `json:"clusterMetrics"`
}
type MetricsInfo struct {
    AppsPending                 int             `json:"appsPending"`
    AppsRunning                 int             `json:"appsRunning"`
    AllocatedMB                 int64           `json:"allocatedMB"`
    AllocatedVirtualCores       int             `json:"allocatedVirtualCores"`
}

func GetClusterInfo(path string) error {
    now_t := time.Now()
    now_str := now_t.Format("2006-01-02 15:04:05")
	request := gorequest.New()
    _, body, errs := request.Get("http://10.142.97.4:8088/ws/v1/cluster/metrics").End()
    if len(errs) != 0 {
        fmt.Errorf("http get error:%v", errs)
    }
    clusterinfo := new(RMMetrics)
    err := json.Unmarshal([]byte(body), clusterinfo)
    if err != nil {
        fmt.Printf("aaa:%s\n", err)
        fmt.Printf("aaa:%s\n", body)
    }
    metrics_arr := make([]string, 0)
    metrics_arr = append(metrics_arr, now_str)
    metrics_arr = append(metrics_arr, strconv.Itoa(clusterinfo.MetricsObj.AppsPending))
    metrics_arr = append(metrics_arr, strconv.Itoa(clusterinfo.MetricsObj.AppsRunning))
    metrics_arr = append(metrics_arr, strconv.FormatInt(clusterinfo.MetricsObj.AllocatedMB, 10))
    metrics_arr = append(metrics_arr, strconv.Itoa(clusterinfo.MetricsObj.AllocatedVirtualCores))
    WriteCsvFile(path, metrics_arr)

    return nil
}
func WriteCsvFile(file_name string, values []string) error {
    var f *os.File

    f,err := os.OpenFile(file_name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0660)
    if(err != nil){
        panic(err)
    }
    if err != nil {
        panic(err)
    }
	defer f.Close()
    w := csv.NewWriter(f)

    w.Write(values)
    w.Flush()
    return nil
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main(){
    path := flag.String("path", "metrics.csv", "file path")
    flag.Parse()
    GetClusterInfo(*path)
}
