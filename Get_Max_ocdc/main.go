package main

import (
    "fmt"
    "os"
    "flag"
    //"bufio"
    "encoding/csv"
    "strings"
    "time"
	"sort"
    "strconv"
    "io"
    /*"github.com/parnurzeal/gorequest"
    "encoding/json"*/
    //"os/exec"
    //"sync"
    //"runtime"
)

type AppProperty struct {
    Queue                       string
    User                        string
    Id                          string
    Name                        string
    StartedTime                 string
    FinishedTime                string
    RealCores                   int
    VcoreSeconds                string
    MemorySeconds               string
    ElapsedTime                 string
    State                       string
}

type AppInfos []AppProperty
func (c AppInfos) Len() int {
	return len(c)
}
func (c AppInfos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c AppInfos) Less(i, j int) bool {
	return c[i].RealCores > c[j].RealCores
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
    for _, info := range infos {
        if target_num > 0 {
            request_number += 1
            if request_number > target_num {
                break;
            }
        }
        app_ret_arr := make([]string, 0)
        app_ret_arr = append(app_ret_arr, info.Queue, info.User, info.Id, info.Name, info.StartedTime, info.FinishedTime, strconv.Itoa(info.RealCores), info.VcoreSeconds, info.MemorySeconds, info.ElapsedTime, info.State)
        w.Write(app_ret_arr)
        w.Flush()
    }

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

func ParseBigFile(topfile, bigfile string) error {
    f,err := os.Open(bigfile)
    if(err != nil){
        panic(err)
    }
	defer f.Close()
    reader := csv.NewReader(f)

    m_BigApp := make(map[string]AppProperty)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            fmt.Println("finish")
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            return nil
        }
        //fmt.Println(record)
        new_app := new(AppProperty)
        new_app.Queue = record[0]
        new_app.User = record[1]
        new_app.Id = record[2]
        new_app.Name = record[3]
        new_app.StartedTime = record[4]
        new_app.FinishedTime = record[5]
        new_cpu_v, _ := strconv.Atoi(record[6])
        new_app.RealCores = new_cpu_v
        new_app.VcoreSeconds = record[7]
        new_app.MemorySeconds = record[8]
        new_app.ElapsedTime = record[9]
        new_app.State = record[10]
        if old_val, ok := m_BigApp[record[2]]; ok {
            if old_val.RealCores < new_app.RealCores {
                m_BigApp[record[2]] = *new_app
            }
        } else {
            m_BigApp[record[2]] = *new_app
        }
    }

    sorted_big_src := make(AppInfos, 0)
    for _, v := range m_BigApp {
        sorted_big_src = append(sorted_big_src, v)
    }
	if !sort.IsSorted(sorted_big_src) {
		sort.Sort(sorted_big_src)
	}
    bignames := strings.Split(bigfile, ".")
    sorted_src_path := bignames[0] + "_sorted_src.csv"
    WriteAllCsvValues(sorted_big_src, sorted_src_path, -1)

    if topfile != "-" {
        ParseTop(m_BigApp, topfile, bigfile)
    }
    return nil
}
func ParseTop(m_BigApp map[string]AppProperty, topfile, bigfile string) error {
    f,err := os.Open(topfile)
    if(err != nil){
        panic(err)
    }
	defer f.Close()
    reader := csv.NewReader(f)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            fmt.Println("finish")
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            return nil
        }
        topnames := strings.Split(topfile, ".")
        realcore_file := topnames[0] + "_real.csv"
        real_core := make([]string, 0)
        if val, ok := m_BigApp[record[2]]; ok {
            real_core = append(real_core, val.Id, strconv.Itoa(val.RealCores))
            WriteCsvFile(realcore_file, real_core)
            delete(m_BigApp, record[2])
        } else {
            real_core = append(real_core, record[2], "0")
            WriteCsvFile(realcore_file, real_core)
        }
    }
    sorted_big_src := make(AppInfos, 0)
    for _, v := range m_BigApp {
        sorted_big_src = append(sorted_big_src, v)
    }
	if !sort.IsSorted(sorted_big_src) {
		sort.Sort(sorted_big_src)
	}
    bignames := strings.Split(bigfile, ".")
    sorted_src_path := bignames[0] + "_sorted_delete.csv"
    WriteAllCsvValues(sorted_big_src, sorted_src_path, -1)
    return nil
}
func main(){
    topfile := flag.String("topfile", "top50_0605.csv", "path")
    bigfile := flag.String("bigfile", "big_app_2018-06-05.csv", "path")
    flag.Parse()
    //fmt.Println(*queue_name, *log_path)
    ParseBigFile(*topfile, *bigfile)
}
