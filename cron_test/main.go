package main

import (
    "github.com/robfig/cron"
    "log"
	"fmt"
    "time"
    "encoding/json"
    "os/exec"
    "strings"
)

type CronTest struct {
    Name        string
    value       int
}

var i = 0
var j = 0
var c *cron.Cron
func (CronTest)cb_cron1(){
    i++
    log.Println("start--i", i)
    c.Stop()
    c1 := cron.New()
    spec2 := "31 * 9 * * ?"
    var ct CronTest
    c1.AddFunc(spec2, ct.cb_cron2)
    c1.Start()
}
func (CronTest)cb_cron2(){
    j++
    log.Println("start--j", j)
}

type TestJob struct {
    Value   int
    Name string
}
func (s *TestJob) Run() {
    fmt.Print(s.Name)
}
func test1() {
    c = cron.New()
    spec1 := "1 * 9 * * ?"
    var ct CronTest
    c.AddFunc(spec1, ct.cb_cron1)
    c.Start()
    select{} //阻塞主线程不退出
}
func test_job(c *cron.Cron) {
    spec1 := "1 7 18 * * 1"
    spec2 := "1 10 18 * * 1"
    c.Start()
    time.Sleep(5)
    //var t TestJob
    t1 := new(TestJob)
    t2 := new(TestJob)
    t1.Value = 1
    t1.Name = "czq 11111"
    c.AddJob(spec1, t1)
    t2.Name = "czq 22222"
    t2.Value = 2
    c.AddJob(spec2, t2)
    //select{} //阻塞主线程不退出
}
func test_time_in(period_type string, start_week, start_hour, start_minute, stop_week, stop_hour, stop_minute int) bool {
    t := time.Now()
    current_week := int(t.Weekday())
    current_hour := t.Hour()
    current_minute := t.Minute()
    current_week = 4
    current_hour = 14
    current_minute = 50
    //fmt.Printf("current:%d, :%d, :%d\n", current_week, current_hour, current_minute)
    //fmt.Printf("start:  %d, :%d, :%d\n", start_week, start_hour, start_minute)
    //fmt.Printf("stop   :%d, :%d, :%d\n", stop_week, stop_hour, stop_minute)

    if period_type == "week"{
        if current_week > start_week && current_week < stop_week {
            return true
        } else if current_week < start_week || current_week > stop_week {
            return false
        } else {
            if stop_week == start_week {
                if current_week == start_week {
                    if start_hour == stop_hour {
                        if current_minute > start_minute && current_minute < stop_minute {
                            return true
                        } else {
                            return false
                        }
                    } else if start_hour < stop_hour {
                        if current_hour > start_hour && current_hour < stop_hour {
                            return true
                        } else if current_hour == start_hour {
                            if current_minute > start_minute {
                                return true
                            } else {
                                return false
                            }
                        } else if current_hour == stop_hour {
                            if current_minute < stop_minute {
                                return true
                            } else {
                                return false
                            }
                        } else {
                            return false
                        }
                    } else {
                        return false
                    }
                } else {
                    return false
                }
            } else if start_week < stop_week {
                if current_week == start_week {
                    if current_hour > start_hour {
                        return true
                    } else if current_hour == start_hour {
                        if current_minute > start_minute {
                            return true
                        } else {
                            return false
                        }
                    } else {
                        return false
                    }
                } else if current_week == stop_week {
                    if current_hour < stop_hour {
                        return true
                    } else if current_hour == stop_hour {
                        if current_minute < stop_minute {
                            return true
                        } else {
                            return false
                        }
                    } else {
                        return false
                    }
                } else {
                    return false
                }
            } else {
                return false
            }
        }
    } else {
        if start_hour == stop_hour {
            if current_minute > start_minute && current_minute < stop_minute {
                return true
            } else {
                return false
            }
        } else if start_hour < stop_hour {
            if current_hour > start_hour && current_hour < stop_hour {
                return true
            } else if current_hour == start_hour {
                if current_minute > start_minute {
                    return true
                } else {
                    return false
                }
            } else if current_hour == stop_hour {
                if current_minute < stop_minute {
                    return true
                } else {
                    return false
                }
            } else {
                return false
            }
        } else {
            return false
        }
    }
    return false
}
func test_time_check(period_type string, start_week, start_hour, start_minute, stop_week, stop_hour, stop_minute int) bool {
    fmt.Printf("start:  %d, :%d, :%d\n", start_week, start_hour, start_minute)
    fmt.Printf("stop   :%d, :%d, :%d\n", stop_week, stop_hour, stop_minute)
    if period_type == "week"{
        if start_week < stop_week {
            return true
        } else if start_week > stop_week {
            return false
        } else {
            if start_hour < stop_hour {
                return true
            } else if start_hour > stop_hour {
                return false
            } else {
                if start_minute < stop_minute {
                    return true
                } else {
                    return false
                }
            }
        }
    } else {
        if start_hour < stop_hour {
            return true
        } else if start_hour > stop_hour {
            return false
        } else {
            if start_minute < stop_minute {
                return true
            } else {
                return false
            }
        }
    }
}

func test_time_check2(period_type string, start_week, start_hour, start_minute, stop_week, stop_hour, stop_minute int) bool {
    start_value := 0
    stop_value := 0
    if period_type == "week"{
        start_value = start_week * 24 * 40 + start_hour * 60 + start_minute
        stop_value = stop_week * 24 * 40 + stop_hour * 60 + stop_minute
        return stop_value > start_value
    } else {
        start_value = start_hour * 60 + start_minute
        stop_value =  stop_hour * 60 + stop_minute
    }
    return stop_value > start_value
}
func test_Overlop(period_type string, start_week1, start_hour1, start_minute1, stop_week1, stop_hour1, stop_minute1, start_week2, start_hour2, start_minute2, stop_week2, stop_hour2, stop_minute2 int) bool {
    fmt.Printf("%d,%d,%d-----%d,%d,%d\n", start_week1, start_hour1, start_minute1, stop_week1, stop_hour1, stop_minute1)
    fmt.Printf("%d,%d,%d-----%d,%d,%d\n", start_week2, start_hour2, start_minute2, stop_week2, stop_hour2, stop_minute2)
    if period_type == "week"{
        start1_value := start_week1 * 24 * 40 + start_hour1 * 60 + start_minute1
        stop1_value := stop_week1 * 24 * 40 + stop_hour1 * 60 + stop_minute1
        start2_value := start_week2 * 24 * 40 + start_hour2 * 60 + start_minute2
        stop2_value := stop_week2 * 24 * 40 + stop_hour2 * 60 + stop_minute2
        return isOverlap_1(start1_value, stop1_value, start2_value, stop2_value)
    } else {
        start1_value := start_hour1 * 60 + start_minute1
        stop1_value :=  stop_hour1 * 60 + stop_minute1
        start2_value := start_hour2 * 60 + start_minute2
        stop2_value :=  stop_hour2 * 60 + stop_minute2
        return isOverlap_1(start1_value, stop1_value, start2_value, stop2_value)
    }
}

func isOverlap_1(start1, end1, start2, end2 int) bool {
    if end1 < start2 || start1 > end2 {
        return false
    }
    return true
}

func isInInterval(current, start, end int) bool {
    if current < start || current > end {
        return false
    }
    return true
}

func test_time_in1(period_type string, start_week, start_hour, start_minute, stop_week, stop_hour, stop_minute int) bool {
    current_week := 4
    current_hour := 14
    current_minute := 50
    if period_type == "week"{
        start_value := start_week * 24 * 40 + start_hour * 60 + start_minute
        stop_value := stop_week * 24 * 40 + stop_hour * 60 + stop_minute
        current_value := current_week * 24 * 40 + current_hour * 60 + current_minute
        return isInInterval(current_value, start_value, stop_value)
    } else {
        start_value := start_hour * 60 + start_minute
        stop_value :=  stop_hour * 60 + stop_minute
        current_value := current_hour * 60 + current_minute
        return isInInterval(current_value, start_value, stop_value)
    }
}

func test_time_cases(){
//	test1()
   fmt.Printf("%t\n", test_time_in("week", 3, 12, 22, 5, 10, 11))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 22, 5, 14, 11))
   fmt.Printf("%t\n", test_time_in("week", 4, 11, 1, 3, 10, 11))
   fmt.Printf("%t\n", test_time_in("week", 4, 11, 31, 4, 11, 11))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 31, 4, 14, 11))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 31, 4, 15, 11))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 1, 4, 14, 36))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 31, 4, 14, 59))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 31, 4, 14, 31))
   fmt.Printf("%t\n", test_time_in("week", 4, 14, 50, 4, 14, 40))
   fmt.Printf("%t\n", test_time_in("week", 4, 13, 50, 4, 12, 40))

   fmt.Println("--------------------")
   fmt.Printf("%t\n", test_time_in1("week", 3, 12, 22, 5, 10, 11))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 22, 5, 14, 11))
   fmt.Printf("%t\n", test_time_in1("week", 4, 11, 1, 3, 10, 11))
   fmt.Printf("%t\n", test_time_in1("week", 4, 11, 31, 4, 11, 11))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 31, 4, 14, 11))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 31, 4, 15, 11))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 1, 4, 14, 36))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 31, 4, 14, 59))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 31, 4, 14, 31))
   fmt.Printf("%t\n", test_time_in1("week", 4, 14, 50, 4, 14, 40))
   fmt.Printf("%t\n", test_time_in1("week", 4, 13, 50, 4, 12, 40))

   fmt.Println("--------------------test_time_check")

   fmt.Printf("%t\n", test_time_check("week", 4, 14, 31, 5, 14, 59))
   fmt.Printf("%t\n", test_time_check("week", 4, 14, 31, 4, 15, 59))
   fmt.Printf("%t\n", test_time_check("week", 4, 14, 31, 4, 14, 59))
   fmt.Printf("%t\n", test_time_check("week", 5, 14, 31, 4, 14, 59))
   fmt.Printf("%t\n", test_time_check("week", 4, 14, 31, 4, 13, 59))
   fmt.Println("--------------------test_time_check")
   fmt.Printf("%t\n", test_time_check2("week", 4, 14, 31, 5, 14, 59))
   fmt.Printf("%t\n", test_time_check2("week", 4, 14, 31, 4, 15, 59))
   fmt.Printf("%t\n", test_time_check2("week", 4, 14, 31, 4, 14, 59))
   fmt.Printf("%t\n", test_time_check2("week", 5, 14, 31, 4, 14, 59))
   fmt.Printf("%t\n", test_time_check2("week", 4, 14, 31, 4, 13, 59))
   fmt.Println("--------------------")
   fmt.Printf("%t\n", test_Overlop("week", 4, 14, 31, 4, 15, 59, 5, 13, 30, 5, 14, 30))
   fmt.Printf("%t\n", test_Overlop("week", 4, 14, 31, 4, 15, 59, 4, 13, 30, 5, 14, 30))
   fmt.Printf("%t\n", test_Overlop("week", 4, 14, 31, 4, 15, 59, 3, 13, 30, 3, 14, 30))
   fmt.Printf("%t\n", test_Overlop("week", 4, 14, 31, 4, 15, 59, 3, 13, 30, 5, 14, 30))
   fmt.Printf("%t\n", test_Overlop("week", 4, 14, 31, 4, 15, 59, 3, 13, 30, 5, 14, 30))
}
func GetReplicas(service_type, service_name string) (int32, error){
    checkcmd := "oc get statefulsets spark-worker-mj2 -o json"
    out, err := exec.Command("bash", "-c", checkcmd).CombinedOutput()
    fmt.Printf(string(out))
    if err != nil && strings.Contains(string(out), "not found") {
        return 0, fmt.Errorf("RealResetService patch out:%s", string(out))
    }
    var sp interface{}
    err = json.Unmarshal(out, &sp)
    if err != nil {
        return 0, fmt.Errorf("GetReplicas :%s", err.Error())
    }
    var ret error
    var replicas int32
    if m_sp, ok := sp.(map[string]interface{}); ok {
        if status_v, exist := m_sp["status"]; exist {
            if m_status, ok := status_v.(map[string]interface{}); ok {
                if replicas_i, exist := m_status["replicas"]; exist {
                    if replicas_int, ok := replicas_i.(int32); ok {
                        replicas = replicas_int
                    }
                } else {
                    ret = fmt.Errorf("replicas exist error")
                }
            } else {
                ret = fmt.Errorf("status type error")
            }
        } else {
            ret = fmt.Errorf("status exist error")
        }
    }
    return replicas, ret
}
func GetModuleReadyReplicas(module_type, module_name, project_name string) (int32, error){
    checkcmd := "oc get " + module_type
    checkcmd += " " + module_name
    checkcmd += " -n " + project_name
    checkcmd += " -o json"
    replicas_name := ""
    lower_module_type := strings.ToLower(module_type)
    if lower_module_type == "statefulsets" ||
        lower_module_type == "rc" ||
        lower_module_type == "replicationcontrollers" {
            replicas_name = "replicas"
    } else {
            replicas_name = "readyReplicas"
    }
    //} else if lower_module_type == "rc" || lower_module_type == "replicationcontrollers"
    out, err := exec.Command("bash", "-c", checkcmd).CombinedOutput()
    if err != nil && strings.Contains(string(out), "not found") {
        return -1, fmt.Errorf("GetModuleReadyReplicas patch out:%s", string(out))
    }
    var sp interface{}
    err = json.Unmarshal(out, &sp)
    if err != nil {
        return -1, fmt.Errorf("GetReplicas :%s", err.Error())
    }
    var ret error
    var replicas int32
    if m_sp, ok := sp.(map[string]interface{}); ok {
        if status_v, exist := m_sp["status"]; exist {
            if m_status, ok := status_v.(map[string]interface{}); ok {
                if replicas_i, exist := m_status[replicas_name]; exist {
                    if replicas_int, ok := replicas_i.(float64); ok {
                        replicas = int32(replicas_int)
                    } else {
                        fmt.Printf("... :%d", int(replicas_int))
                        ret = fmt.Errorf("replicas type error")
                    }
                } else {
                    ret = fmt.Errorf("replicas exist error:%s", replicas_name)
                }
            } else {
                ret = fmt.Errorf("status type error")
            }
        } else {
            ret = fmt.Errorf("status exist error")
        }
    }
    return replicas, ret
}
func GetDesireReplicas(module_type, module_name, project_name string) (int32, error){
    checkcmd := "oc get " + module_type
    checkcmd += " " + module_name
    checkcmd += " -n " + project_name
    checkcmd += " -o json"
    out, err := exec.Command("bash", "-c", checkcmd).CombinedOutput()
    if err != nil && strings.Contains(string(out), "not found") {
        return -1, fmt.Errorf("GetModuleReadyReplicas patch out:%s", string(out))
    }
    var sp interface{}
    err = json.Unmarshal(out, &sp)
    if err != nil {
        return -1, fmt.Errorf("GetReplicas :%s", err.Error())
    }
    var ret error
    var replicas int32
    if m_sp, ok := sp.(map[string]interface{}); ok {
        if status_v, exist := m_sp["spec"]; exist {
            if m_status, ok := status_v.(map[string]interface{}); ok {
                if replicas_i, exist := m_status["replicas"]; exist {
                    if replicas_int, ok := replicas_i.(float64); ok {
                        fmt.Printf("aa:%f\n", replicas_int)
                        replicas = int32(replicas_int)
                    } else {
                        fmt.Print(replicas_i)
                        fmt.Print("\n")
                    }
                } else {
                    ret = fmt.Errorf("replicas exist error")
                }
            } else {
                ret = fmt.Errorf("status type error")
            }
        } else {
            ret = fmt.Errorf("status exist error")
        }
    }
    return replicas, ret
}
func GetContainerName(service_type, service_name string) (name string, ret error) {
    checkcmd := "oc get statefulsets spark-worker-mj2 -o json"
    out, err := exec.Command("bash", "-c", checkcmd).CombinedOutput()
    fmt.Printf(string(out))
    if err != nil && strings.Contains(string(out), "not found") {
        return "", fmt.Errorf("RealResetService patch out:%s", string(out))
    }
    var sp interface{}
    err = json.Unmarshal(out, &sp)
    if err != nil {
        return "", err
    }
    if m_sp, ok := sp.(map[string]interface{}); ok {
        if spec_v, exist := m_sp["spec"]; exist {
            if spec_m, ok := spec_v.(map[string]interface{}); ok {
                if template_i, exist := spec_m["template"]; exist {
                    if template_m, ok := template_i.(map[string]interface{}); ok {
                        if interspec_i, exist := template_m["spec"]; exist {
                            if interspec_m, ok := interspec_i.(map[string]interface{}); ok {
                                if containers_i, exist := interspec_m["containers"]; exist {
                                    if containers_a, ok := containers_i.([]interface{}); ok {
                                        if len(containers_a) != 1 {
                                            fmt.Printf("arr len is : %d", len(containers_a))
                                        }
                                        if contaner_m, ok := containers_a[0].(map[string]interface{}); ok {
                                            if name_v, exist := contaner_m["name"]; exist {
                                                if name_s, ok := name_v.(string); ok {
                                                    name = name_s
                                                }
                                            } else {
                                                ret = fmt.Errorf("GetContainerName name exist error")
                                            }
                                        } else {
                                            ret = fmt.Errorf("GetContainerName name type error")
                                        }
                                    } else {
                                        ret = fmt.Errorf("GetContainerName container type error")
                                    }
                                } else {
                                    ret = fmt.Errorf("GetContainerName container exist error")
                                }
                            } else {
                                ret = fmt.Errorf("GetContainerName interspec type error")
                            }
                        } else {
                            ret = fmt.Errorf("GetContainerName interspec exist error")
                        }
                    } else {
                        ret = fmt.Errorf("GetContainerName template type error")
                    }
                } else {
                    ret = fmt.Errorf("GetContainerName template exist error")
                }
            } else {
                ret = fmt.Errorf("GetContainerName spec type error")
            }
        } else {
            ret = fmt.Errorf("GetContainerName spec exist error")
        }
    }
    return name, ret
}
func main() {
    //aa := fmt.Sprintf("{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"spark-worker-mj2\",\"resources\":{\"limits\":{\"cpu\":\"%d\",\"memory\":\"%dG\"}}}]}}}}", 4, 8)
    //fmt.Println(aa)
    //GetReplicas()
    //name, _ := GetContainerName("", "")
    //fmt.Print(name)
    //r, err := GetDesireReplicas("statefulsets", "spark-worker-mj2", "czq-project")
    //r, err := GetModuleReadyReplicas("rc", "nodejs-mongodb-example-15", "openshift")
    /*r, err := GetModuleReadyReplicas("deployment", "accp-d", "czq-project")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("ready:%d", r)
    r, err = GetDesireReplicas("deploy", "accp-d", "czq-project")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("desire:%d", r)*/
/*    c := cron.New()
    go test_job(c)
    time.Sleep(12000 * time.Millisecond)
    c.Stop()
    select{} //阻塞主线程不退出*/
    test1()
}
