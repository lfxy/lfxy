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
    "os/exec"
    "sync"
    "runtime"
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
	return c[i].VcoreSeconds > c[j].VcoreSeconds
	//return c[i].AllocatedVCores > c[j].AllocatedVCores
	//return c[i].RealCores > c[j].RealCores
}

func GetQueueInfo(new_path, all_path, ip_tail string) error {
    now_t := time.Now()
    now_str := now_t.Format("2006-01-02 15:04:05")
    new_stamp := now_t.Unix()
	request := gorequest.New()
    str_url := fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?user=ocdc", ip_tail)
    _, body, errs := request.Get(str_url).End()
    if len(errs) != 0 {
        fmt.Errorf("http get error:%v", errs)
    }
    //fmt.Printf("body:%d\n", len(body))

    /*
    fileName := "c.json"
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
    }*/
    //fmt.Printf("body:%s\n", body)
    //fmt.Printf("now_str:%s\n", now_str)
    //fmt.Printf("now_stamp:%d\n", new_stamp)
    rminfo := new(RMApps)
    err := json.Unmarshal([]byte(body), rminfo)
    if err != nil {
        fmt.Println(err)
    }
	if !sort.IsSorted(rminfo.AppsObj.AppObj) {
		sort.Sort(rminfo.AppsObj.AppObj)
	}

    new_apps := make([]string, 0)
    new_apps = append(new_apps, now_str)
    all_apps := make([]string, 0)
    all_apps = append(all_apps, now_str)
    for _, app := range rminfo.AppsObj.AppObj {
        all_apps = append(all_apps, app.Id)
        app_start_time := app.StartedTime / 1000
        time_interval := new_stamp - app_start_time
        if time_interval <= 6 && strings.Contains(app.Name, "test_p_ta_cockpit_kpi_detail_month_01"){
            fmt.Printf("new:%s\n", app.Id)
            new_apps = append(new_apps, app.Id)
        }
    }
    if len(new_apps) > 1 {
        fmt.Printf("new:%v\n", new_apps)
        WriteCsvFile(new_path, new_apps)
    }
    //WriteCsvFile(all_path, all_apps)

    return nil
}
func KillApps(state, queue_name, ip_tail string) error {
    str_url := ""
    if state == "all" {
        str_url = fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?queue=%s", ip_tail, queue_name)
    } else {
        str_url = fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?states=%s&queue=%s", ip_tail, state, queue_name)
    }
	request := gorequest.New()
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

    cmd1 := "yarn application -kill "
    count := 1
    fmt.Println("total:%d", len(rminfo.AppsObj.AppObj))
    var wg sync.WaitGroup
    runtime.GOMAXPROCS(8)
    for index, app := range rminfo.AppsObj.AppObj {
        wg.Add(1)
        go func(id string) {
            fmt.Println(count)
            count += 1
            r_cmd := cmd1 + id
            fmt.Printf("%s\n", r_cmd)
            cmd := exec.Command("sh", "-c", r_cmd)
            out, _ := cmd.Output()
            fmt.Printf("out:%s\n", out)
            wg.Done()
            }(app.Id)
        if index % 100 == 0 {
            time.Sleep(60 * time.Second)
        }
    }
    wg.Wait()
    return nil
}
func getCurrentPath() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
		log.Fatal(err)
   }
   return strings.Replace(dir, "\\", "/", -1)
}
func GetAppResourcesByAverage(queue_name, task_status, starttime_begin, starttime_end string, target_num int, ip_tail string) error {
    timeLayout := "2006-01-02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    t_begein, _ := time.ParseInLocation(timeLayout, starttime_begin, loc)
    t_end, _ := time.ParseInLocation(timeLayout, starttime_end, loc)
    ts_begin := t_begein.UnixNano()
    ts_end := t_end.UnixNano()
	request := gorequest.New()
    str_url := ""
    if task_status == "all" {
        str_url = fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?startedTimeBegin=%s&startedTimeEnd=%s", ip_tail, strconv.FormatInt(ts_begin/1000000, 10), strconv.FormatInt(ts_end/1000000, 10))
    } else if task_status == "RUNNING" {
        str_url = fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?state=RUNNING", ip_tail)
    } else {
        str_url = fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?state=%s&startedTimeBegin=%s&startedTimeEnd=%s", ip_tail, task_status, strconv.FormatInt(ts_begin/1000000, 10), strconv.FormatInt(ts_end/1000000, 10))
    }
    if queue_name != "all" {
        str_url += "&user="
        str_url += queue_name
    }
    fmt.Println(queue_name)
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
    var ta AppInfosAverage
    //var ta AppInfosBy
    for _, app := range rminfo.AppsObj.AppObj {
        if task_status == "RUNNING" {
            app.RealCores = app.AllocatedVCores
        } else {
            tt := app.ElapsedTime / 1000
            if tt == 0 {
                fmt.Println("0:", app.State)
                continue
            }
            app.RealCores = (int)(app.VcoreSeconds / tt)
        }
        ta = append(ta, app)
    }
	if !sort.IsSorted(ta) {
		sort.Sort(ta)
	}

    now_t := time.Now().Format("2006-01-02")
	current_path := getCurrentPath()
    big_app_path := current_path + "/big_app_"
    big_app_path += now_t
    big_app_path += ".csv"
    b_big_app, _ := PathExists(big_app_path)
    if !b_big_app {
        app_ret_arr := make([]string, 0)
        app_ret_arr = append(app_ret_arr, "queue", "user", "id", "name", "started_time", "finished_time", "cpu_average", "cpu_seconds", "memory_seconds", "time_elapsed", "finalStatus")
        WriteCsvFile(big_app_path, app_ret_arr)
    }
    fmt.Println("queue,user,id,name,started_time,finished_time,cpu_average,cpu_seconds,memory_seconds,time_elapsed,finalStatus")
    count := 0
    for _, app := range ta {
        if queue_name == "ocdc" || (!strings.Contains(app.Queue, "ocdc") && !strings.Contains(app.Queue, "wzfw") && !strings.Contains(app.Queue, "hxxt_fp") ) {
            app_start_time := time.Unix(app.StartedTime / 1000, 0).Format(timeLayout)
            app_finished_time := time.Unix(app.FinishedTime / 1000, 0).Format(timeLayout)
            fmt.Printf("%s,%s,%s,%s,%s,%s,%d,%d,%d,%d,%s\n", app.Queue, app.User, app.Id, app.Name, app_start_time, app_finished_time, app.RealCores, app.VcoreSeconds, app.MemorySeconds, app.ElapsedTime, app.State)
            count += 1
            if target_num > 0 && count >= target_num {
                break
            }

            if task_status == "RUNNING" && app.RealCores > 1000 {
                app_ret_arr := make([]string, 0)
                app_ret_arr = append(app_ret_arr, app.Queue, app.User, app.Id, app.Name, app_start_time, app_finished_time, strconv.Itoa(app.RealCores), strconv.FormatInt(app.VcoreSeconds, 10), strconv.FormatInt(app.MemorySeconds, 10), strconv.FormatInt(app.ElapsedTime, 10), app.State)
                WriteCsvFile(big_app_path, app_ret_arr)
            }
        }
    }

    return nil
}
/*func GetAppResourcesByTotal(task_status, starttime_begin, starttime_end string) error {
    timeLayout := "2006-01-02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    t_begein, _ := time.ParseInLocation(timeLayout, starttime_begin, loc)
    t_end, _ := time.ParseInLocation(timeLayout, starttime_end, loc)
    ts_begin := t_begein.UnixNano()
    ts_end := t_end.UnixNano()
	request := gorequest.New()
    str_url := ""
    //str_url := fmt.Sprintf("http://10.142.97.3:8088/ws/v1/cluster/apps?state=FINISHED&applicationTypes=MAPREDUCE&startedTimeBegin=%s&startedTimeEnd=%s", strconv.FormatInt(ts_begin/1000000, 10), strconv.FormatInt(ts_end/1000000, 10))
    if task_status == "all" {
        str_url = fmt.Sprintf("http://10.142.97.3:8088/ws/v1/cluster/apps?startedTimeBegin=%s&startedTimeEnd=%s", strconv.FormatInt(ts_begin/1000000, 10), strconv.FormatInt(ts_end/1000000, 10))
    } else {
        str_url = fmt.Sprintf("http://10.142.97.3:8088/ws/v1/cluster/apps?state=%s&startedTimeBegin=%s&startedTimeEnd=%s", task_status, strconv.FormatInt(ts_begin/1000000, 10), strconv.FormatInt(ts_end/1000000, 10))

    }
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

    fmt.Println("queue,id,name,started_time,cpu_used,memory_used,time_elapsed,finalStatus")
    count := 0
    for _, app := range rminfo.AppsObj.AppObj {
        if !strings.Contains(app.Queue, "root.ocdc") {
            app_start_time := time.Unix(app.StartedTime / 1000, 0).Format(timeLayout)
            app_finished_time := time.Unix(app.FinishedTime / 1000, 0).Format(timeLayout)
            fmt.Printf("%s,%s,%s,%s,%s,%s,%d,%d,%d,%s\n", app.Queue, app.User, app.Id, app.Name, app_start_time, app_finished_time, app.VcoreSeconds, app.MemorySeconds, app.ElapsedTime, app.State)
            count += 1
            if count >= 50 {
                break
            }
        }
    }

    return nil
}*/
func GetRunningAppResources(task_type, starttime_begin, starttime_end, ip_tail string) error {
    timeLayout := "2006-01-02 15:04:05"
	request := gorequest.New()
    str_url := fmt.Sprintf("http://10.142.97.%s:8088/ws/v1/cluster/apps?state=RUNNING", ip_tail)
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

    fmt.Println("queue,id,name,started_time,cpu,memory")
    now_t := time.Now().Format("2006-01-02 15:04:05")
    count := 0
    for _, app := range rminfo.AppsObj.AppObj {
        if !strings.Contains(app.Queue, "root.ocdc") {
            app_start_time := time.Unix(app.StartedTime / 1000, 0).Format(timeLayout)
            fmt.Printf("%s,%s,%s,%s,%s,%d,%d\n", now_t,app.Queue, app.Id, app.Name, app_start_time, app.AllocatedVCores, app.AllocatedMB)
            count += 1
            if count >= 50 {
                break
            }
        }
    }

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

//   ./ocdc_apps -o=top -name=all -state=RUNNING
//   ./ocdc_apps -o=top -name=all -state=FINISHED -begin='2018-06-05 17:00:00' -end='2018-06-06 09:30:00'
//  ./ocdc_apps -o=kill -name=ocdc -state=ACCEPTED

func main(){
    //new_path := flag.String("newpath", "new_app.csv", "file path")
    //all_path := flag.String("allpath", "all_app.csv", "file path")
    operation := flag.String("o", "top", "operation select")
    begin := flag.String("begin", "2018-05-29 09:20:00", "file path")
    end := flag.String("end", "2018-05-31 09:20:00", "file path")
    task_state := flag.String("state", "all", "file path")
    queue_name := flag.String("name", "hjpt", "file path")
    target_num := flag.Int("num", 50, "path")
    ip_tail := flag.String("ip", "4", "path")
    //queue_state := flag.String("state", "ACCEPTED", "file path")
    flag.Parse()
//    GetQueueInfo(*new_path, *all_path)
    if *operation == "top" {
        GetAppResourcesByAverage(*queue_name, *task_state, *begin, *end, *target_num, *ip_tail)
    } else if *operation == "kill" {
        KillApps(*task_state, *queue_name, *ip_tail)
    } else {
        fmt.Println("error operation beside top and kill:", *operation)
    }
}
