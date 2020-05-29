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
    XMLName         xml.Name        `xml:"allocations"`
    RootObj         RootIns         `xml:"pool"`
}
type RootIns struct {
    XMLName     xml.Name        `xml:"pool"`
    PoolObjs    []PoolIns       `xml:"pool"`
}
type PoolIns struct {
    XMLName     xml.Name                `xml:"pool"`
    PoolName    string                  `xml:"name,attr"`
    Min_v       string                  `xml:"minResources"`
    Max_v       string                  `xml:"maxResources"`
    Acl_v       string                  `xml:"aclSubmitApps"`
    PoolSeconds []PoolIns               `xml:"pool"`
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
    for _, pool := range v.RootObj.PoolObjs {
        if pool.PoolSeconds != nil {
            //fmt.Println("====================")
            for _, pool2 := range pool.PoolSeconds {
                mem_int, cpu_int, err := ParseMemCpu(pool2.Min_v)
                if err != nil {
                    fmt.Println(err.Error())
                    return
                }
                mem_int2, cpu_int2, err := ParseMemCpu(pool2.Max_v)
                if err != nil {
                    fmt.Println(err.Error())
                    return
                }
                fmt.Printf("%s=%d=%d=%d=%d=root.normal_queues.bigplatform.%s.%s\n", strings.TrimSuffix(pool2.Acl_v, ","), mem_int, cpu_int, mem_int2, cpu_int2, pool.PoolName, pool2.PoolName)
            }
            //fmt.Println("++++++++++===========")
            continue
        }
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
        fmt.Printf("%s=%d=%d=%d=%d=root.normal_queues.bigplatform.%s\n", strings.TrimSuffix(pool.Acl_v, ","), mem_int, cpu_int, mem_int2, cpu_int2, pool.PoolName)
    }
}
func test(fileName string){
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
    }
    var v XML_test
    err = xml.Unmarshal(data, &v)
    if err != nil {
        fmt.Println(err)
    }
    var min_mem_total int64
    var min_cpu_total int
    var max_mem_total int64
    var max_cpu_total int

    var max_min_jt_mem int64
    var max_min_jt_cpu int
    var max_max_jt_mem int64
    var max_max_jt_cpu int
    for _, pool := range v.RootObj.PoolObjs {
        mem_int, cpu_int, err := ParseMemCpu(pool.Min_v)
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        min_mem_total += mem_int
        min_cpu_total += cpu_int

        mem_int2, cpu_int2, err := ParseMemCpu(pool.Max_v)
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        max_mem_total += mem_int2
        max_cpu_total += cpu_int2

        /*if !strings.HasPrefix(pool.PoolName, "jt_") {
            continue
        }*/
        fmt.Printf("%s,%d,%d,%d,%d\n", pool.PoolName, mem_int, cpu_int, mem_int2, cpu_int2)
        //fmt.Println(pool.Min_v)
        if mem_int > max_min_jt_mem {
            max_min_jt_mem = mem_int
        }
        if cpu_int > max_min_jt_cpu {
            max_min_jt_cpu = cpu_int
        }


        if mem_int2 > max_max_jt_mem {
            max_max_jt_mem = mem_int2
        }
        if cpu_int2 > max_max_jt_cpu {
            max_max_jt_cpu = cpu_int2
        }
    }
    fmt.Printf("min_mem_total,%d, min_cpu_total,%d\n", min_mem_total, min_cpu_total)
    fmt.Printf("max_mem_total,%d, max_cpu_total,%d\n", max_mem_total, max_cpu_total)
    //fmt.Printf("max_min_jt_mem:%d, max_min_jt_cpu:%d\n", max_min_jt_mem, max_min_jt_cpu)
    //fmt.Printf("max_max_jt_mem:%d, max_max_jt_cpu:%d\n", max_max_jt_mem, max_max_jt_cpu)
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

type ConfigurationIns struct {
    XMLName     xml.Name        `xml:"configuration"`
    PropertyObjs    []PropertyIns       `xml:"property"`
}
type PropertyIns struct {
    XMLName     xml.Name        `xml:"property"`
    Name    string          `xml:"name"`
    Value    string          `xml:"value"`
    Source       string        `xml:"source"`
}
func ParseNarmal(fileName string){
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
    }
    var v ConfigurationIns
    err = xml.Unmarshal(data, &v)
    if err != nil {
        fmt.Println(err)
    }
    for _, p := range v.PropertyObjs {
        fmt.Printf("%s=%s\n", p.Name, p.Value)
    }
}
func main(){
    path := flag.String("path", "bigplatform.xml", "file path")
    flag.Parse()
    test2(*path)
   //ParseNarmal(*path)
}
