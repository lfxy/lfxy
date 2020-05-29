package main

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    "flag"
//    "os"
)

type XML_test struct {
    XMLName         xml.Name            `xml:"allocations"`
    RootFirsts      PoolFirst           `xml:"pool"`
}
type PoolFirst struct {
    XMLName         xml.Name                `xml:"pool"`
    PoolSeconds     []PoolSecond            `xml:"pool"`
    PoolLevel       int
}
type PoolSecond struct {
    XMLName         xml.Name                `xml:"pool"`
    PoolName        string                  `xml:"name,attr"`
    Min_v           string                  `xml:"minResources"`
    Max_v           string                  `xml:"maxResources"`
    Acl_v           string                  `xml:"aclSubmitApps"`
    Acladm_v           string                  `xml:"aclAdministerApps"`
    MaxRunningApps_v           int                  `xml:"maxRunningApps"`
    Weight_v           float32                  `xml:"weight"`
    PoolSeconds     []PoolSecond            `xml:"pool"`
    PoolLevel       int
}

func test2(fileName string){
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
    }
    var v XML_test
    err = xml.Unmarshal(data, &v)
    if err != nil {
        fmt.Println(err)
    }
    v.RootFirsts.PoolLevel = 2
    for _, pool := range v.RootFirsts.PoolSeconds {
        if pool.Min_v != "" && pool.Max_v != "" {
            mem_int, cpu_int, err := ParseMemCpu(pool.Min_v)
            if err != nil {
                fmt.Println(err.Error())
                return
            }
            mem_int2, cpu_int2, err := ParseMemCpu(pool.Max_v)
            if err != nil {
                fmt.Println(err.Error())
                return
            }
            fmt.Printf("%d=%s=%d=%d=%d=%d=root.%s=%s\n", v.RootFirsts.PoolLevel, strings.TrimSuffix(pool.Acl_v, ","), mem_int, cpu_int, mem_int2, cpu_int2, pool.PoolName, strings.TrimSuffix(pool.Acladm_v, ","))
        }

        if pool.PoolSeconds != nil {
            for _, pool2 := range pool.PoolSeconds {
                ParseSecondQueue(pool2, "root." + pool.PoolName, v.RootFirsts.PoolLevel)
            }
        }
    }
}

func ParseSecondQueue(pool2 PoolSecond, fathername string, poolLevel int) error {
    //fmt.Println("====================")
    poolLevel++
    pool2.PoolLevel = poolLevel
    mem_int, cpu_int, err := ParseMemCpu(pool2.Min_v)
    if err != nil {
        fmt.Println(err.Error())
        return fmt.Errorf("ParseSecondQueue error min")
    }
    mem_int2, cpu_int2, err := ParseMemCpu(pool2.Max_v)
    if err != nil {
        fmt.Println(err.Error())
        return fmt.Errorf("ParseSecondQueue error max")
    }
    currentPoolName := fathername + "." + pool2.PoolName
    //fmt.Printf("line%d=%s=%d=%d=%d=%d=%s\n", poolLevel, strings.TrimSuffix(pool2.Acl_v, ","), mem_int, cpu_int, mem_int2, cpu_int2, currentPoolName)
    fmt.Printf("line%d=%s=%d=%d=%d=%d=%d=%.1f=%s=%s\n", poolLevel, strings.TrimSuffix(pool2.Acl_v, ","), mem_int, cpu_int, mem_int2, cpu_int2, pool2.MaxRunningApps_v, pool2.Weight_v , currentPoolName, strings.TrimSuffix(pool2.Acladm_v, ","))
    //fmt.Println("++++++++++===========")

    if pool2.PoolSeconds != nil {
        for _, pool3 := range pool2.PoolSeconds {
            ParseSecondQueue(pool3, currentPoolName, poolLevel)
        }
    }
    return nil
}

func ParseMemCpu(str string) (int64, int, error) {
    keys := strings.Split(str, ",")
    if len(keys) != 2 {
        fmt.Printf("error!")
        return -1, -1, fmt.Errorf("ParseMemCpu error1:%s", str)
    }
    cpu_str := strings.TrimPrefix(keys[1], " ")
    mems := strings.Split(keys[0], " ")
    cpus := strings.Split(cpu_str, " ")

    mem_int, err1 := strconv.ParseInt(mems[0], 10, 64)
    cpu_int,err2:=strconv.Atoi(cpus[0])
    if err1 != nil || err2 != nil {
        fmt.Printf("error !!!!!")
        return -1, -1, fmt.Errorf("ParseMemCpu error2")
    }

    return mem_int, cpu_int, nil
}

func main(){
    path := flag.String("path", "fair-scheduler.xml", "1_line=2_user=3_minMem=4_min_cpu=5_max_mem=6_max_cpu=7_max_apps=8_weight=9_queue=10_admuser")
    flag.Parse()
    //fmt.Printf("1_line=2_user=3_minMem=4_min_cpu=5_max_mem=6_max_cpu=7_max_apps=8_weight=9_queue=10_admuser\n")
    test2(*path)
    //./xml2 | grep -v  bigplatform | grep line5 | awk -F '=' '{print $7","$3","$4","$5","$6}'
}
