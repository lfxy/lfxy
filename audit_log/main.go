package main

import (
    "fmt"
    "os"
    "flag"
    "bufio"
    "encoding/csv"
    "strings"
    "time"
	"sort"
    "strconv"
    /*"github.com/parnurzeal/gorequest"
    "encoding/json"*/
    //"os/exec"
    //"sync"
    //"runtime"
)

type AppProperty struct {
    ApplicationId                   string
    Status                          string
    SubmitTime                      string
    StartTime                       string
    FinishedTime                    string
    AllocateTime                    string
    TotalTime                       string
    AllocateDuration                int64
    TotalDuration                   int64
}

type AppInfos []AppProperty
func (c AppInfos) Len() int {
	return len(c)
}
func (c AppInfos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c AppInfos) Less(i, j int) bool {
	return c[i].AllocateDuration > c[j].AllocateDuration
	//return c[i].TotalDuration > c[j].TotalDuration
	//return c[i].SubmitTime < c[j].SubmitTime
}

func SubTime(begin, end string) (string, int64) {
    timeLayout := "2006/01/02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    t_begin, _ := time.ParseInLocation(timeLayout, begin, loc)
    t_end, _ := time.ParseInLocation(timeLayout, end, loc)

    d_ret := t_end.Sub(t_begin)
    //fmt.Println(d_ret.String())
    p_hours := (int)(d_ret.Hours())
    p_minutes := (int)(d_ret.Minutes()) - p_hours * 60
    p_seds := (int)(d_ret.Seconds()) - (int)(d_ret.Minutes()) * 60
    //s_ret := fmt.Sprintf("0%d:%d:%d", (int)(d_ret.Hours()), (int)(d_ret.Minutes()), (int)(d_ret.Seconds()))
    s_ret := fmt.Sprintf("0%d:%d:%d", p_hours, p_minutes, p_seds)

    return s_ret, (int64)(d_ret.Seconds())
}
func WriteCsvFile(file_name string, values []string) error {
    f,err := os.OpenFile(file_name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0660)
    if(err != nil){
        panic(err)
    }
	defer f.Close()
    w := csv.NewWriter(f)

    w.Write(values)
    w.Flush()
    return nil
}
func WriteAllCsvValues(infos AppInfos, file_name string, target_num int) error {
    f,err := os.OpenFile(file_name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0660)
    if(err != nil){
        panic(err)
    }
	defer f.Close()
    w := csv.NewWriter(f)
    request_number := 0
    var total_allocate int64
    var total_total int64 
    for _, info := range infos {
        request_number += 1
        if target_num > 0 && request_number > target_num {
                break;
        }
        app_ret_arr := make([]string, 0)
        app_ret_arr = append(app_ret_arr, info.ApplicationId, info.Status, info.SubmitTime, info.StartTime, info.FinishedTime, info.AllocateTime, info.TotalTime, strconv.FormatInt(info.AllocateDuration, 10), strconv.FormatInt(info.TotalDuration, 10))
        w.Write(app_ret_arr)
        w.Flush()

        total_allocate += info.AllocateDuration
        total_total += info.TotalDuration
    }
    average_allocate := total_allocate / (int64)(request_number)
    average_total := total_total / int64(request_number)
    app_ret_arr := make([]string, 0)
    app_ret_arr = append(app_ret_arr, "", "", "", "", "", "", "", strconv.FormatInt(average_allocate, 10), strconv.FormatInt(average_total, 10))
    w.Write(app_ret_arr)
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

func ParseAuditLog(queue_name, log_path, des_path, request_begin, request_end string, target_num int) error {
    m_apps := make(map[string]*AppProperty, 0)
    f,err := os.Open(log_path)
    if(err != nil){
        panic(err)
    }
	defer f.Close()
    br := bufio.NewReader(f)
    for {
        str, err := br.ReadString('\n') //每次读取一行
        if err!= nil {
            break
        }
        if strings.Contains(str, queue_name) {
            //fmt.Printf(str)
            if strings.Contains(str, "Submit Application Request") {
                str_sub := strings.Split(str, "\t")
                sub_time := strings.Split(str_sub[0], "|")
                app_id_arr := strings.Split(str_sub[5], "=")
                app_id := strings.TrimSuffix(app_id_arr[1], "\n")
                app_obj := new(AppProperty)
                app_obj.ApplicationId = app_id
                app_obj.SubmitTime = "20" + sub_time[0]
                m_apps[app_obj.ApplicationId] = app_obj
            } else if strings.Contains(str, "Register App Master") {
                str_sub := strings.Split(str, "\t")
                app_id_arr := strings.Split(str_sub[5], "=")
                app_id := strings.TrimSuffix(app_id_arr[1], "\n")
                if app_line, ok := m_apps[app_id]; ok {
                    start_time := strings.Split(str_sub[0], "|")
                    app_line.StartTime = "20" + start_time[0]
                    s_allocate, i_allocate := SubTime(app_line.SubmitTime, app_line.StartTime)
                    app_line.AllocateTime = s_allocate
                    app_line.AllocateDuration = i_allocate
                    m_apps[app_id] = app_line
                } else {
                    fmt.Printf("error !!, Register App Master not exist id:%s\n%s\n", app_id, str)
                }
            } else if strings.Contains(str, "Application Finished - Succeeded") {
                str_sub := strings.Split(str, "\t")
                app_id_arr := strings.Split(str_sub[4], "=")
                app_id := strings.TrimSuffix(app_id_arr[1], "\n")
                if app_line, ok := m_apps[app_id]; ok {
                    finished_time := strings.Split(str_sub[0], "|")
                    states := strings.Split(str_sub[3], "=")
                    app_line.Status = states[1]
                    app_line.FinishedTime = "20" + finished_time[0]

                    s_total, i_total := SubTime(app_line.SubmitTime, app_line.FinishedTime)
                    app_line.TotalTime = s_total
                    app_line.TotalDuration = i_total
                    m_apps[app_id] = app_line
                } else {
                    fmt.Printf("error !!, Application Finished not exist id:%s\n%s", app_id, str)
                }
            }
        }
    }

    //fmt.Printf("application_id,status,submit_time,start_time,finished_time,allocate_time,total_time,allocate_seconds,total_seconds")
    app_ret_arr := make([]string, 0)
    app_ret_arr = append(app_ret_arr, "application_id", "status", "submit_time", "start_time", "finished_time", "allocate_time", "total_time", "allocate_seconds", "total_seconds")
    WriteCsvFile(des_path, app_ret_arr)

    var infos AppInfos
    for _, v := range m_apps {
        if v.TotalTime == "" {
            continue
        }
        if request_begin != "" {
            if v.SubmitTime < request_begin || v.SubmitTime > request_end {
                continue
            }
        }
        infos = append(infos, *v)

    }
	if !sort.IsSorted(infos) {
		sort.Sort(infos)
	}

    WriteAllCsvValues(infos, des_path, target_num)

/*    for k, v := range m_apps {
        if v.TotalTime == "" {
            continue
        }
        app_ret_arr := make([]string, 0)
        //fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n", k, v.Status, v.SubmitTime, v.StartTime, v.FinishedTime, v.AllocateTime, v.TotalTime)
        app_ret_arr = append(app_ret_arr, k, v.Status, v.SubmitTime, v.StartTime, v.FinishedTime, v.AllocateTime, v.TotalTime)
        WriteCsvFile(des_path, app_ret_arr)
    }*/
    return nil
}
// ./audit_log -path=./yarn-rm.audit.log.2018-05-26 -begin="2018/05/26 02:00:00" -end="2018/05/26 08:00:00" -num=2000 -des=./des_26_3.csv
// ./audit_log -path=./yarn-rm.audit.log -begin="2018/06/01 00:00:00" -end="2018/06/01 08:00:00" -num=2000 -des=./des_01_3.csv
func main(){
    queue_name := flag.String("name", "ocdc", "name")
    log_path := flag.String("path", "./yarn-rm.audit.log", "path")
    des_path := flag.String("des", "./des.csv", "path")
    request_begin := flag.String("begin", "", "path")
    request_end := flag.String("end", "", "path")
    target_num := flag.Int("num", -1, "path")
    flag.Parse()
    //fmt.Println(*queue_name, *log_path)
    ParseAuditLog(*queue_name, *log_path, *des_path, *request_begin, *request_end, *target_num)
}
