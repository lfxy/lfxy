package main

import (
    "fmt"
    "os/exec"
    "flag"
    "strings"
    "strconv"
	"sort"
    "time"
    //"influxdbclient"
    //"github.com/parnurzeal/gorequest"
    //"os"
    //"log"
    //"path/filepath"
    //"encoding/json"
    //"encoding/csv"
    //"sync"
    //"runtime"
)
type LoadAvg struct {
    HostName            string
    Value               float64
}

type LoadAvgs []LoadAvg
func (c LoadAvgs) Len() int {
	return len(c)
}
func (c LoadAvgs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c LoadAvgs) Less(i, j int) bool {
	return c[i].Value > c[j].Value
}

    //cmd1 := "pssh -i -h nmlist.txt \"uptime;hostname;\" | grep -v \"SUCCESS\" | grep -v \"known hosts\""
func Get_Loads(file_name, key_path string, target_value float64, i_min int) (int, error) {
    ret_num := -1
    now_t := time.Now().Format("2006-01-02 15:04:05")
    cmd1 := fmt.Sprintf("/usr/bin/parallel-ssh -i -h %s -x \"-i %s\" \"uptime;hostname;\" | grep -v \"SUCCESS\" | grep -v \"known hosts\"", file_name, key_path)
    //cmd1 := fmt.Sprintf("ansible all -i %s -m shell -a \"uptime;hostname;\" | grep -v \"SUCCESS\" | grep -v \"known hosts\"", file_name)
    fmt.Println(cmd1)
    cmd := exec.Command("bash", "-c", cmd1)
    out, err := cmd.Output()
    str_out := string(out)
    if err != nil || !strings.Contains(str_out, "load average") {
        fmt.Println(err)
        return ret_num, fmt.Errorf("err!")
    }
    //fmt.Printf("out:%s\n", out)
    lines := strings.Split(str_out, "\n")
    /*m1 := make(map[string]float64, 0)
    m5 := make(map[string]float64, 0)
    m15 := make(map[string]float64, 0)*/
    var m1 LoadAvgs
    var m5 LoadAvgs
    var m15 LoadAvgs
    for i, line := range lines {
        fmt.Println("1:" + line)
        if strings.Contains(line, "load average") {
            load_avgs := strings.Split(line, ",")
            load_str := ""
            avg_5 := ""
            avg_15 := ""
            if strings.Contains(load_avgs[3], "load") {
                load_str = load_avgs[3]
                avg_5 = load_avgs[4]
                avg_15 = load_avgs[5]
            } else if strings.Contains(load_avgs[2], "load") {
                load_str = load_avgs[2]
                avg_5 = load_avgs[3]
                avg_15 = load_avgs[4]
            } else {
                fmt.Println("err load_f:" + line)
                continue
            }
            load_f := strings.Split(load_str, ":")
            v1, err1 := strconv.ParseFloat(strings.TrimPrefix(load_f[1], " "), 64)
            v5, err5 := strconv.ParseFloat(strings.TrimPrefix(avg_5, " "), 64)
            v15, err15 := strconv.ParseFloat(strings.TrimPrefix(avg_15, " "), 64)
            if err1 != nil || err5 != nil || err15 != nil {
                fmt.Errorf("err!")
                break
            }
            if v1 > target_value {
                //m1[lines[i + 1]] = v1
                var load1 LoadAvg
                load1.HostName = lines[i + 1]
                load1.Value = v1
                m1 = append(m1, load1)
                //fmt.Printf("%f       %s\n", v1, lines[i + 1])
                if v5 > target_value {
                    //m5[lines[i + 1]] = v5
                    var load5 LoadAvg
                    load5.HostName = lines[i + 1]
                    load5.Value = v5
                    m5 = append(m5, load5)
                }
                if v15 > target_value {
                    //m15[lines[i + 1]] = v15
                    var load15 LoadAvg
                    load15.HostName = lines[i + 1]
                    load15.Value = v15
                    m15 = append(m15, load15)
                }
            }
        }
    }
	if !sort.IsSorted(m1) {
		sort.Sort(m1)
	}
	if !sort.IsSorted(m5) {
		sort.Sort(m5)
	}
	if !sort.IsSorted(m15) {
		sort.Sort(m15)
	}
    if i_min == 1 {
        fmt.Printf("load average 1 minute=%s=%d\n", now_t, len(m1))
        for _, v1 := range m1 {
            fmt.Printf("%f       %s\n", v1.Value, v1.HostName)
        }
        ret_num = len(m1)
    } else if i_min == 5 {
        ret_num = len(m5)
        fmt.Printf("load average 1 minute=%s=%d\n", now_t, len(m1))
        fmt.Printf("load average 5 minute=%s=%d\n", now_t, len(m5))
        for _, v5 := range m5 {
            fmt.Printf("%f       %s\n", v5.Value, v5.HostName)
        }
    } else if i_min == 15 {
        //fmt.Printf("load average 1 minute=%s=%d\n", now_t, len(m1))
        //fmt.Printf("load average 15 minute=%s=%d\n", now_t,len(m15))
        fmt.Printf("%d\n", len(m15))
        ret_num = len(m15)
        /*for _, v15 := range m15 {
            fmt.Printf("%f       %s\n", v15.Value, v15.HostName)
        }*/
    } else {
        ret_num = len(m15)
        fmt.Printf("load average 1 minute=%s=%d\n", now_t ,len(m1))
        for _, v1 := range m1 {
            fmt.Printf("%f       %s\n", v1.Value, v1.HostName)
        }
        fmt.Printf("load average 5 minute=%s=%d\n", now_t, len(m5))
        for _, v5 := range m5 {
            fmt.Printf("%f       %s\n", v5.Value, v5.HostName)
        }
        fmt.Printf("load average 15 minute=%s=%d\n", now_t, len(m15))
        for _, v15 := range m15 {
            fmt.Printf("%f       %s\n", v15.Value, v15.HostName)
        }
    }
    return ret_num, nil
}
func Get_Centos7(file_name string, target_value float64, i_min int) error {
    //cmd1 := fmt.Sprintf("pssh -i -h %s \"uptime;hostname;\" | grep -v \"SUCCESS\" | grep -v \"known hosts\"", file_name)
    cmd1 := "pssh -i -h all_ip.txt \"cat /etc/redhat-release;hostname;\" | grep -v SUCCESS | grep -v \"known hosts\""
    cmd := exec.Command("sh", "-c", cmd1)
    out, _ := cmd.Output()
    lines := strings.Split(string(out), "\n")
    m1 := make([]string, 0)
    for i, line := range lines {
        //fmt.Println(line)
        if strings.Contains(line, "release 7") {
            //fmt.Println("aa:"+lines[i + 1])
            m1 = append(m1, lines[i + 1])
        }
    }
    for _, hostn := range m1 {
        fmt.Println(hostn)
    }
    return nil
}
func main(){
    //file_name := flag.String("name", "/home/op/tmp/czq2/get_load/nmlist.txt", "operation select")
    file_name := flag.String("name", "/home/devops/czq/hosts.txt", "hosts file")
    target_value := flag.Float64("value", 32, "load threashold")
    i_min := flag.Int("min", 15, "which time load")
    key_path := flag.String("key", "/home/devops/.ssh/ladevops_key", "ssh key file")
    //queue_state := flag.String("state", "ACCEPTED", "file path")
    flag.Parse()
    ret_num, err := Get_Loads(*file_name, *key_path, *target_value, *i_min)
    if err != nil {
       fmt.Printf("error:%s\n", err)
       return
    }
    var influxCli *DbClient = new(DbClient)
    influxCli.Cli = DbConfig()
	influxCli.WritePoint("yarn_crontab_loadaverage", "realtime.yarn", ret_num)
    //Get_Centos7(*file_name, *target_value, *i_min)
}
