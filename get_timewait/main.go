package main

import (
    "fmt"
    "os/exec"
    "flag"
    "strings"
    "strconv"
	//"sort"
    //"time"
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
func Get_Close_Wait_Num(file_name, key_path string) (map[string]interface {}, error) {
    mhostnum := make(map[string]interface{})
    cmd1 := fmt.Sprintf("/usr/bin/parallel-ssh -i -h %s -x \"-i %s\" \"hostname;sudo netstat -nat | grep 9083 | grep CLOSE_WAIT | wc -l;\" | grep -v \"SUCCESS\" ", file_name, key_path)
    cmd := exec.Command("bash", "-c", cmd1)
    out, err := cmd.Output()
    str_out := string(out)
    if err != nil || !strings.Contains(str_out, "ip-") {
        return mhostnum, fmt.Errorf("err:%v", err)
    }
    //fmt.Printf("out:%s\n", out)
    lines := strings.Split(str_out, "\n")
    for i, line := range lines {
        if strings.Contains(line, "ip-") {
            v1, err1 := strconv.Atoi(lines[i + 1])
            if err1 != nil {
                fmt.Printf("parse num error:%v", err1)
                continue
            }
            mhostnum[line] = v1
        }
    }
    return mhostnum, nil
}
func main(){
    //file_name := flag.String("name", "/home/op/tmp/czq2/get_load/nmlist.txt", "operation select")
    file_name := flag.String("name", "/home/devops/czq/get_load/metastore_hosts.txt", "hosts file")
    key_path := flag.String("key", "/home/devops/.ssh/ladevops_key", "ssh key file")
    tag := flag.String("tag", "cdh2.metastore", "tag name")
    flag.Parse()
    fields, err := Get_Close_Wait_Num(*file_name, *key_path)
    if err != nil {
       fmt.Printf("error:%s\n", err)
       return
    }
    var influxCli *DbClient = new(DbClient)
    influxCli.Cli = DbConfig()
    influxCli.WritePoint("metastore_crontab_closewait", *tag, fields)
}
