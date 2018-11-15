package main

import (
    "fmt"
    "github.com/parnurzeal/gorequest"
    "os"
    "io"
    "encoding/json"
	"sort"
    "encoding/csv"
    "flag"
    "time"
    "strconv"
    "strings"
    "log"
    "path/filepath"
)

type RMResponse struct {
    SchedulerObj         RMInfo          `json:"scheduler"`
}
type RMInfo struct {
    InfoObj             InfoMessage          `json:"schedulerInfo"`
}

type InfoMessage struct {
    Type        string          `json:"type"`
    RootQueue   RootInfo        `json:"rootQueue"`
}
type RootInfo struct {
    QueueInfo
    ChildQueues            QueueInfos         `json:"childQueues"`
}
type QueueInfos []QueueInfo
type QueueInfo struct {
    Type                    string          `json:"type,omitempty"`
    MaxApps                 int             `json:"maxApps"`
    MinResources            ResourceInfo    `json:"minResources"`
    MaxResources            ResourceInfo    `json:"maxResources"`
    UsedResources           ResourceInfo    `json:"usedResources"`
    SteadyFairResources     ResourceInfo    `json:"steadyFairResources"`
    FairResources           ResourceInfo    `json:"fairResources"`
    ClusterResources        ResourceInfo    `json:"clusterResources"`
    QueueName               string          `json:"queueName"`
    SchedulingPolicy        string          `json:"schedulingPolicy"`
    NumPendingApps          int             `json:"numPendingApps,omitempty"`
    NumActiveApps           int             `json:"numActiveApps,omitempty"`
}
type ResourceInfo struct {
    Memory          int64           `json:"memory"`
    Vcores          int64           `json:"vCores"`
}

func (c QueueInfos) Len() int {
	return len(c)
}
func (c QueueInfos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c QueueInfos) Less(i, j int) bool {
	return c[i].UsedResources.Memory > c[j].UsedResources.Memory
}

func GetRmInfo(cpu_path, mem_path, min_max_path string) error {
	request := gorequest.New()
    _, body, errs := request.Get("http://10.142.97.3:8088/ws/v1/cluster/scheduler").End()
    if len(errs) != 0 {
        fmt.Errorf("http get error:%v", errs)
        _, body, errs = request.Get("http://10.142.97.4:8088/ws/v1/cluster/scheduler").End()
        if len(errs) != 0 {
            return fmt.Errorf("http request error:%v", errs)
        }
    }
    //fmt.Printf("body:%d\n", len(body))

    /*
    fileName := "c.json"
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
    }*/
    rminfo := new(RMResponse)
    err := json.Unmarshal([]byte(body), rminfo)
    if err != nil {
        fmt.Println(err)
    }

    queue_names := make([]string, 0)
    queue_names = append(queue_names, "")

    queue_used_memory := make([]string, 0)
    queue_used_cpu := make([]string, 0)
    queue_pending_app := make([]string, 0)
    queue_running_app := make([]string, 0)

    now_str := time.Now().Format("2006-01-02 15:04:05")
    queue_used_memory = append(queue_used_memory, now_str)
    queue_used_cpu = append(queue_used_cpu, now_str)
    queue_pending_app = append(queue_pending_app, now_str)
    queue_running_app = append(queue_running_app, now_str)


    queue_min_mem := make([]string, 0)
    queue_max_mem := make([]string, 0)
    queue_min_cpu := make([]string, 0)
    queue_max_cpu := make([]string, 0)

    queue_min_mem = append(queue_min_mem, "min_memory")
    queue_min_cpu = append(queue_min_cpu, "min_cpu")
    queue_max_mem = append(queue_max_mem, "max_memory")
    queue_max_cpu = append(queue_max_cpu, "max_cpu")

    var mem_total_used int64
    var cpu_total_used int64
    min_max_exist, _ := PathExists(min_max_path)
    cpu_exist, _ := PathExists(cpu_path)
    mem_exist, _ := PathExists(mem_path)
    if cpu_exist && mem_exist {
        fmt.Println("x ..............")
        m_queues_info := make(map[string]QueueInfo, 0)
        for _, queueinfo := range rminfo.SchedulerObj.InfoObj.RootQueue.ChildQueues {
            m_queues_info[queueinfo.QueueName] = queueinfo
        }
		mem_file, err := os.Open(min_max_path)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		defer mem_file.Close()
		reader := csv.NewReader(mem_file)
		record, err := reader.Read()
        if err != nil || err == io.EOF {
            fmt.Println("Error:", err)
            return err
        }
        //fmt.Println(record)
        for _, queue_name := range record {
            if !strings.HasPrefix(queue_name, "root") {
                continue
            }
            if queue_info_t, ok := m_queues_info[queue_name]; ok {
                queue_used_memory = append(queue_used_memory, strconv.FormatInt(queue_info_t.UsedResources.Memory, 10))
                queue_used_cpu = append(queue_used_cpu, strconv.FormatInt(queue_info_t.UsedResources.Vcores, 10))
                mem_total_used += queue_info_t.UsedResources.Memory
                cpu_total_used += queue_info_t.UsedResources.Vcores
                queue_pending_app = append(queue_pending_app, strconv.Itoa(queue_info_t.NumPendingApps))
                queue_running_app = append(queue_running_app, strconv.Itoa(queue_info_t.NumActiveApps))
                if queue_info_t.NumPendingApps > 0 {
                    fmt.Printf("%s  %s pending apps is :%d\n", now_str, queue_info_t.QueueName, queue_info_t.NumPendingApps)
                }

                queue_min_mem = append(queue_min_mem, strconv.FormatInt(queue_info_t.MinResources.Memory, 10))
                queue_max_mem = append(queue_max_mem, strconv.FormatInt(queue_info_t.MaxResources.Memory, 10))
                queue_min_cpu = append(queue_min_cpu, strconv.FormatInt(queue_info_t.MinResources.Vcores, 10))
                queue_max_cpu = append(queue_max_cpu, strconv.FormatInt(queue_info_t.MaxResources.Vcores, 10))
                delete(m_queues_info, queue_name)
            } else {
                queue_used_memory = append(queue_used_memory, "0")
                queue_used_cpu = append(queue_used_cpu, "0")
                queue_pending_app = append(queue_pending_app, "0")
                queue_running_app = append(queue_running_app, "0")

                queue_min_mem = append(queue_min_mem, "0")
                queue_max_mem = append(queue_max_mem, "0")
                queue_min_cpu = append(queue_min_cpu, "0")
                queue_max_cpu = append(queue_max_cpu, "0")
            }
        }
        for _, v := range m_queues_info {
            fmt.Println("new queue ..............:", v.QueueName)
            record = append(record, v.QueueName)
            queue_used_memory = append(queue_used_memory, strconv.FormatInt(v.UsedResources.Memory, 10))
            queue_used_cpu = append(queue_used_cpu, strconv.FormatInt(v.UsedResources.Vcores, 10))

            queue_min_mem = append(queue_min_mem, strconv.FormatInt(v.MinResources.Memory, 10))
            queue_max_mem = append(queue_max_mem, strconv.FormatInt(v.MaxResources.Memory, 10))
            queue_min_cpu = append(queue_min_cpu, strconv.FormatInt(v.MinResources.Vcores, 10))
            queue_max_cpu = append(queue_max_cpu, strconv.FormatInt(v.MaxResources.Vcores, 10))
            //SeekCsvFile(mem_path, record)
            //SeekCsvFile(cpu_path, record)
            //WriteCsvFile(mem_path, record)
            //WriteCsvFile(cpu_path, record)
        }
        if len(m_queues_info) > 0 {
            os.Remove(min_max_path)
            WriteCsvFile(min_max_path, record)
            WriteCsvFile(min_max_path, queue_min_mem)
            WriteCsvFile(min_max_path, queue_max_mem)
            WriteCsvFile(min_max_path, queue_min_cpu)
            WriteCsvFile(min_max_path, queue_max_cpu)
        }
    } else {
        fmt.Println("First ..............")
        if !sort.IsSorted(rminfo.SchedulerObj.InfoObj.RootQueue.ChildQueues) {
            sort.Sort(rminfo.SchedulerObj.InfoObj.RootQueue.ChildQueues)
        }
        for _, queueinfo := range rminfo.SchedulerObj.InfoObj.RootQueue.ChildQueues {
            queue_names = append(queue_names, queueinfo.QueueName)
            queue_used_memory = append(queue_used_memory, strconv.FormatInt(queueinfo.UsedResources.Memory, 10))
            queue_used_cpu = append(queue_used_cpu, strconv.FormatInt(queueinfo.UsedResources.Vcores, 10))
            mem_total_used += queueinfo.UsedResources.Memory
            cpu_total_used += queueinfo.UsedResources.Vcores

            if !min_max_exist {
                queue_min_mem = append(queue_min_mem, strconv.FormatInt(queueinfo.MinResources.Memory, 10))
                queue_max_mem = append(queue_max_mem, strconv.FormatInt(queueinfo.MaxResources.Memory, 10))
                queue_min_cpu = append(queue_min_cpu, strconv.FormatInt(queueinfo.MinResources.Vcores, 10))
                queue_max_cpu = append(queue_max_cpu, strconv.FormatInt(queueinfo.MaxResources.Vcores, 10))
            }
            queue_pending_app = append(queue_pending_app, strconv.Itoa(queueinfo.NumPendingApps))
            queue_running_app = append(queue_running_app, strconv.Itoa(queueinfo.NumActiveApps))
            if queueinfo.NumPendingApps > 0 {
                fmt.Printf("%s  %s pending apps is :%d\n", now_str, queueinfo.QueueName, queueinfo.NumPendingApps)
            }
        }
        if !min_max_exist {
            WriteCsvFile(min_max_path, queue_names)
            WriteCsvFile(min_max_path, queue_min_mem)
            WriteCsvFile(min_max_path, queue_max_mem)
            WriteCsvFile(min_max_path, queue_min_cpu)
            WriteCsvFile(min_max_path, queue_max_cpu)
        }
    }
    fmt.Printf("%s cpu total:%d, mem total:%d\n", now_str, cpu_total_used, mem_total_used)

	current_path := getCurrentPath()
    pending_path := current_path + "/pending_app.csv"
    running_path := current_path + "/running_app.csv"
    if b_exist, _ := PathExists(pending_path); !b_exist{
        WriteCsvFile(pending_path, queue_names)
    }
    if b_exist, _ := PathExists(running_path); !b_exist{
        WriteCsvFile(running_path, queue_names)
    }
    if b_exist, _ := PathExists(mem_path); !b_exist{
        WriteCsvFile(mem_path, queue_names)
    }
    if b_exist, _ := PathExists(cpu_path); !b_exist{
        WriteCsvFile(cpu_path, queue_names)
    }
    WriteCsvFile(mem_path, queue_used_memory)
    WriteCsvFile(cpu_path, queue_used_cpu)
    WriteCsvFile(pending_path, queue_pending_app)
    WriteCsvFile(running_path, queue_running_app)
    return nil
}

func getCurrentPath() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
		log.Fatal(err)
   }
   return strings.Replace(dir, "\\", "/", -1)
}
func SeekCsvFile(file_name string, values []string) error {
    var f *os.File

    f,err := os.OpenFile(file_name, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        panic(err)
    }
	defer f.Close()
    f.Seek(0,os.SEEK_SET)
    w := csv.NewWriter(f)

    w.Write(values)
    w.Flush()
    return nil
}
func WriteCsvFile(file_name string, values []string) error {
    var f *os.File

    f,err := os.OpenFile(file_name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0660)
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

func main(){
    cpu_path := flag.String("cpu_path", "cpu.csv", "cpu file path")
    mem_path := flag.String("mem_path", "mem.csv", "mem file path")
    min_max_path := flag.String("min_max_path", "min_max.csv", "min and max file path")
    flag.Parse()
    GetRmInfo(*cpu_path, *mem_path, *min_max_path)
}
