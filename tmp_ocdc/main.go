package main

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
    "os"
    "log"
    "path/filepath"
    "encoding/json"
	"sort"
    "encoding/csv"
    "flag"
    "time"
    "strings"
    "strconv"
    /*"os/exec"
    "sync"
    "runtime"*/
)

type RMApps struct {
    AppsObj         AppObj          `json:"apps"`
}
type AppObj struct {
    AppObj             AppInfos          `json:"app"`
}
type AppInfos []AppProperty
type AppInfosAverage []AppProperty

type AppProperty struct {
    FinishedTime                int64               `json:"finishedTime"`
    AmContainerLogs             string              `json:"amContainerLogs"`
    TrackingUI                  string              `json:"trackingUI"`
    State                       string              `json:"state"`
    User                        string              `json:"user"`
    Id                          string              `json:"id"`
    ClusterId                   int64               `json:"clusterId"`
    FinalStatus                 string              `json:"finalStatus"`
    AmHostHttpAddress           string              `json:"amHostHttpAddress"`
    Progress                    float32             `json:"progress"`
    Name                        string              `json:"name"`
    StartedTime                 int64               `json:"startedTime"`
    ElapsedTime                 int64               `json:"elapsedTime"`
    Diagnostics                 string              `json:"diagnostics"`
    TrackingUrl                 string              `json:"trackingUrl"`
    Queue                       string              `json:"queue"`
    AllocatedMB                 int                 `json:"allocatedMB"`
    AllocatedVCores             int                 `json:"allocatedVCores"`
    RunningContainers           int                 `json:"runningContainers"`
    MemorySeconds               int64               `json:"memorySeconds"`
    VcoreSeconds                int64               `json:"vcoreSeconds"`
    RealCores                   int               `json:"realcore,omitempty"`
}

func (c AppInfosAverage) Len() int {
	return len(c)
}
func (c AppInfosAverage) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c AppInfosAverage) Less(i, j int) bool {
	//return c[i].VcoreSeconds > c[j].VcoreSeconds
	//return c[i].AllocatedVCores > c[j].AllocatedVCores
	return c[i].RealCores > c[j].RealCores
}
func (c AppInfos) Len() int {
	return len(c)
}
func (c AppInfos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c AppInfos) Less(i, j int) bool {
	return c[i].ElapsedTime > c[j].ElapsedTime
	//return c[i].AllocatedVCores > c[j].AllocatedVCores
	//return c[i].RealCores > c[j].RealCores
}

func getCurrentPath() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
		log.Fatal(err)
   }
   return strings.Replace(dir, "\\", "/", -1)
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

func GetAppResourcesByAverage(user_name, task_status string, target_num int) error {
    timeLayout := "2006-01-02 15:04:05"
    //loc, _ := time.LoadLocation("Local")
	request := gorequest.New()
    //str_url := fmt.Sprintf("http://10.142.97.3:8088/ws/v1/cluster/apps?user=%s&state=%s", user_name, task_status)
    str_url := fmt.Sprintf("http://10.142.97.3:8088/ws/v1/cluster/apps?state=FINISHED")
    //str_url := fmt.Sprintf("http://192.168.11.136:8088/ws/v1/cluster/apps?user=%s&state=%s", user_name, task_status)
    fmt.Println(str_url)
    _, body, errs := request.Get(str_url).End()
    if len(errs) != 0 {
        fmt.Errorf("http get error:%v", errs)
    }
    rminfo := new(RMApps)
    err := json.Unmarshal([]byte(body), rminfo)
    if err != nil {
        fmt.Println(err)
    }
	if !sort.IsSorted(rminfo.AppsObj.AppObj) {
		sort.Sort(rminfo.AppsObj.AppObj)
	}

    now_t := time.Now().Format("2006-01-02")
	current_path := getCurrentPath()
    big_app_path := current_path + "/app_"
    big_app_path += now_t
    big_app_path += ".csv"
    b_big_app, _ := PathExists(big_app_path)
    if !b_big_app {
        app_ret_arr := make([]string, 0)
        app_ret_arr = append(app_ret_arr, "queue", "user", "id", "name", "started_time", "finished_time", "time_elapsed")
        WriteCsvFile(big_app_path, app_ret_arr)
    }
    fmt.Println("queue,user,id,name,started_time,finished_time,time_elapsed")
    for _, app := range rminfo.AppsObj.AppObj {
        app_start_time := time.Unix(app.StartedTime / 1000, 0).Format(timeLayout)
        app_finished_time := time.Unix(app.FinishedTime / 1000, 0).Format(timeLayout)
        app_elapsed_time := app.ElapsedTime / 1000
        //fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", app.Queue, app.User, app.Id, app.Name, app_start_time, app_finished_time, strconv.FormatInt(app_elapsed_time, 10))

        app_ret_arr := make([]string, 0)
        app_ret_arr = append(app_ret_arr, app.Queue, app.User, app.Id, app.Name, app_start_time, app_finished_time, strconv.FormatInt(app_elapsed_time, 10))
        WriteCsvFile(big_app_path, app_ret_arr)
    }

    return nil
}
//   ./ocdc_apps  -name=ocdc -state=FINISHED

func main(){
    task_state := flag.String("state", "FINISHED", "path")
    user_name := flag.String("name", "ocdc", "path")
    target_num := flag.Int("num", -1, "path")
    flag.Parse()
    GetAppResourcesByAverage(*user_name, *task_state, *target_num)
}
