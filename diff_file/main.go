package main

import (
    "fmt"
    "os"
    "flag"
    "bufio"
    //"encoding/csv"
    "strings"
    //"time"
	//"sort"
    //"strconv"
    /*"github.com/parnurzeal/gorequest"
    "encoding/json"*/
    //"os/exec"
    //"sync"
    //"runtime"
)


func ParseAuditLog(src_path, des_path string) error {
    m_apps := make(map[string]int, 0)
    f,err := os.Open(src_path)
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
        str_arr := strings.Split(str, "\t")
        fmt.Printf("s1:%s\n", str)
        m_apps[str_arr[1]] = 1
    }

    for key, _ := range m_apps {
        fmt.Println(key)
    }

    f2,err := os.Open(des_path)
    if(err != nil){
        panic(err)
    }
	defer f2.Close()
    br2 := bufio.NewReader(f2)
    for {
        str2, err := br2.ReadString('\n') //每次读取一行
        if err!= nil {
            break
        }
        str_arr2 := strings.Split(str2, "\t")
        fmt.Printf("s2:%s\n", str2)
		if _, exists := m_apps[str_arr2[1]]; !exists {
            fmt.Println(str2)
        }
    }
    return nil
}
func main(){
    src_path := flag.String("path", "./7.txt", "path")
    des_path := flag.String("des", "./206.txt", "path")
    flag.Parse()
    //fmt.Println(*queue_name, *log_path)
    ParseAuditLog(*src_path, *des_path)
}
