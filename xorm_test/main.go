package main

import (
	"fmt"
    "encoding/json"
    "strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/lfxy/xorm_test/models"
	"github.com/satori/go.uuid"
    "time"
    "os/exec"
    "strings"
)

func CommunicateProjects(project_uuid, service_addr string, b_creating bool)(error){
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    addr_arr := strings.Split(service_addr, ",")
    addr_part := strings.Split(addr_arr[0], ".")
    if len(addr_part) < 3 {
        return fmt.Errorf("CommunicateProjects service_addr error format:%s", service_addr)
    }
    depend_project_name := ""
    project_name_index := 0
    if addr_part[3] == "svc" {
        depend_project_name = addr_part[2]
        project_name_index = 2
    } else if addr_part[2] == "svc" {
        depend_project_name = addr_part[1]
        project_name_index = 1
    } else {
        return fmt.Errorf("CommunicateProjects service_addr error format and need svc:%s", service_addr)
    }
// check other addrs
    for _, addr_obj := range addr_arr {
        addr_part := strings.Split(addr_obj, ".")
        if len(addr_part) < 3 {
            return fmt.Errorf("CommunicateProjects service_addr error format:%s", service_addr)
        }
        if addr_part[project_name_index + 1] != "svc" {
            return fmt.Errorf("CommunicateProjects in check service_addr error format:%s", service_addr)
        }
        if addr_part[project_name_index] != depend_project_name {
            return fmt.Errorf("CommunicateProjects in check cannot depend more than 1 service:%s", service_addr)
        }
    }

    current_proj, err := GetProjectByUuid(project_uuid)
    if err != nil {
        return fmt.Errorf("CommunicateProjects select data project_uuid:%s\n error:%s", project_uuid, err)
    }
    if current_proj.ProjectName == depend_project_name {
        return nil
    }
    depend_proj, err := GetProjectByName(current_proj.UserId, depend_project_name)
    if err != nil {
        return fmt.Errorf("CommunicateProjects select data user_id:%s, project_name:%s, error:%s", current_proj.UserId, depend_project_name, err)
    }


    if b_creating {
        return HandleProjectsCommunicateCreating(current_proj, depend_proj)
    } else {
        return HandleProjectsCommunicateDeleting(current_proj, depend_proj)
    }
    return nil
}

func HandleProjectsCommunicateCreating(current_proj, depend_proj *Projects)(error){
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    current_project_map := make(map[string]int)
    depend_project_map := make(map[string]int)
    if current_proj.CommunicatedProjects != "" {
        err := json.Unmarshal([]byte(current_proj.CommunicatedProjects), &current_project_map)
        if err != nil {
            return fmt.Errorf("CommunicateProjects parse json current_project err:%s", err.Error())
        }
    }
    if depend_proj.CommunicatedProjects != "" {
        err := json.Unmarshal([]byte(depend_proj.CommunicatedProjects), &depend_project_map)
        if err != nil {
            return fmt.Errorf("CommunicateProjects parse json depend_project err:%s", err.Error())
        }
    }
    depend_project_num, exist := current_project_map[depend_proj.ProjectUuid]
    depend_project_num_dup, exist_dup := depend_project_map[current_proj.ProjectUuid]

    if exist_dup != exist || depend_project_num_dup != depend_project_num {
        return fmt.Errorf("CommunicateProjects depend_num does not equal %s:%s\n%s:%s", current_proj.ProjectUuid, current_proj.CommunicatedProjects, depend_proj.ProjectUuid, depend_proj.CommunicatedProjects)
    }

    if !exist || depend_project_num == 0 {
        network_cmd := "oadm pod-network join-projects --to="
        network_cmd += depend_proj.ProjectName
        network_cmd += " "
        network_cmd += current_proj.ProjectName
        fmt.Printf("CommunicateProjects network_cmd:%s", network_cmd)
        /*out, err := exec.Command("bash", "-c", network_cmd).CombinedOutput()
        if err != nil {
            return fmt.Errorf("CommunicateProjects err:%s, out:%s", err.Error(), string(out))
        }*/
        if !exist {
            current_project_map[depend_proj.ProjectUuid] = 1
        } else {
            current_project_map[depend_proj.ProjectUuid] = depend_project_num + 1
        }
    } else {
        current_project_map[depend_proj.ProjectUuid] = depend_project_num + 1
    }
    current_project_str, err := json.Marshal(current_project_map)
    if err != nil {
        return fmt.Errorf("CommunicateProjects marshal current_project err:%s", err.Error())
    }
    UpdateProjectsTable(current_proj.ProjectUuid, "communicated_projects", string(current_project_str))

    if !exist_dup {
        depend_project_map[current_proj.ProjectUuid] = 1
    } else {
        depend_project_map[current_proj.ProjectUuid] = depend_project_num_dup + 1
    }
    depend_project_str, err := json.Marshal(depend_project_map)
    if err != nil {
        return fmt.Errorf("CommunicateProjects marshal depend_project err:%s", err.Error())
    }
    UpdateProjectsTable(depend_proj.ProjectUuid, "communicated_projects", string(depend_project_str))
    return nil
}

func HandleProjectsCommunicateDeleting(current_proj, depend_proj *Projects)(error){
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    current_project_map := make(map[string]int)
    depend_project_map := make(map[string]int)
    if current_proj.CommunicatedProjects != "" {
        err := json.Unmarshal([]byte(current_proj.CommunicatedProjects), &current_project_map)
        if err != nil {
            return fmt.Errorf("HandleProjectsCommunicateDeleting parse json current_project err:%s", err.Error())
        }
    }
    if depend_proj.CommunicatedProjects != "" {
        err := json.Unmarshal([]byte(depend_proj.CommunicatedProjects), &depend_project_map)
        if err != nil {
            return fmt.Errorf("HandleProjectsCommunicateDeleting parse json depend_project err:%s", err.Error())
        }
    }
    depend_project_num, exist := current_project_map[depend_proj.ProjectUuid]
    depend_project_num_dup, exist_dup := depend_project_map[current_proj.ProjectUuid]

    if exist_dup != exist || depend_project_num_dup != depend_project_num {
        return fmt.Errorf("HandleProjectsCommunicateDeleting depend_num does not equal %s:%s\n%s:%s", current_proj.ProjectUuid, current_proj.CommunicatedProjects, depend_proj.ProjectUuid, depend_proj.CommunicatedProjects)
    }

    if !exist || depend_project_num == 0 {
        return fmt.Errorf("HandleProjectsCommunicateDeleting error depend numbers:%s", current_proj.CommunicatedProjects)
    } else {
        current_project_map[depend_proj.ProjectUuid] = depend_project_num - 1
        if current_project_map[depend_proj.ProjectUuid] == 0 {
            network_cmd := "oadm pod-network isolate-projects "
            network_cmd += depend_proj.ProjectName
            network_cmd += " "
            network_cmd += current_proj.ProjectName
            fmt.Printf("HandleProjectsCommunicateDeleting network_cmd:%s", network_cmd)
            /*out, err := exec.Command("bash", "-c", network_cmd).CombinedOutput()
            if err != nil {
                return fmt.Errorf("HandleProjectsCommunicateDeleting err:%s, out:%s", err.Error(), string(out))
            }*/
        }
    }
    current_project_str, err := json.Marshal(current_project_map)
    if err != nil {
        return fmt.Errorf("HandleProjectsCommunicateDeleting marshal current_project err:%s", err.Error())
    }
    UpdateProjectsTable(current_proj.ProjectUuid, "communicated_projects", string(current_project_str))

    depend_project_map[current_proj.ProjectUuid] = depend_project_num_dup - 1
    depend_project_str, err := json.Marshal(depend_project_map)
    if err != nil {
        return fmt.Errorf("HandleProjectsCommunicateDeleting marshal depend_project err:%s", err.Error())
    }
    UpdateProjectsTable(depend_proj.ProjectUuid, "communicated_projects", string(depend_project_str))
    return nil
}
type Projects struct {
    Id                              int        `json:"id,omitempty"`
    ProjectUuid                     string     `json:"project_uuid"`
	ProjectName                     string     `json:"project_name"`
    UserId                          string     `json:"use_id,omitempty"`
	ProjectType                     string     `json:"cluster_type"`
    BillingMode                     string     `json:"billing_mode"`
    CreateTime                      time.Time  `xorm:"created"`
    EndTime                         time.Time  `json:"end_time,omitempty"`
    Status                          string     `json:"status,omitempty"`
    CpuAccumulate                   float64    `json:"cpu_accumulate,omitempty"`
    MemoryAccumulate                float64    `json:"memory_accumulate,omitempty"`
    S3Accumulate                    float64    `json:"s3_accumulate,omitempty"`
    NfsAccumulate                   float64    `json:"nfs_accumulate,omitempty"`
    IfBilling                       bool       `json:"if_biiling" xorm:"Bool"`
    CommunicatedProjects            string     `json:"communicated_projects,omitempty"`
}
func UpdateProjectsTable(project_uuid, column_name, column_value interface{}) (error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    sql := fmt.Sprintf("update `projects` set %s=? where project_uuid=?", column_name)
    _, rerr := engine.Exec(sql, column_value, project_uuid)
    return rerr
}
func GetProjectByUuid(pid string) (*Projects, error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    proj := new(Projects)
    _, err = engine.Where("projects.project_uuid = ?", pid).Get(proj)
    if err != nil {
        return nil, fmt.Errorf("GetProjectByUuid select db error:%s", err)
    }
    return proj, nil
}
func GetProjectByName(user_id, project_name string) (*Projects, error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    proj := new(Projects)
    _, err = engine.Where("projects.user_id = ?", user_id).And("projects.project_name = ?", project_name).Get(proj)
    if err != nil {
        return nil, fmt.Errorf("GetProjectByName select db error:%s", err)
    }
    return proj, nil
}
func test_user() {
	engine, err := xorm.NewEngine("mysql", "root:restfulapi123@(10.209.224.161:10022)/czq")
	if err != nil {
		fmt.Println("error1")
	}
	defer engine.Close()
	engine.ShowSQL(true)
	sb := make([]models.ServiceBase, 0)
	if err = engine.Find(&sb); err != nil {
		fmt.Println("error2, ", err)
		return
	}
	for _, i := range sb {
		fmt.Println("i:\n", i)
	}
	fmt.Println("==================\n")
	//sb = make([]models.ServiceBase, 0)
	var singlesb models.ServiceBase
	engine.Where("service_base.service_id = ?", "abc123").Get(&singlesb)
	fmt.Println(singlesb)
}

type ServiceModule struct {
	models.ModuleBase  `xorm:"extends"`
	models.ServiceBase `xorm:"extends"`
}

func (ServiceModule) TableName() string {
	return "module_base"
}

func UpdateServicePlanTable(service_plan_uuid, column_name, start_time string) (error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    sql := fmt.Sprintf("update `service_plans` set %s=? where service_plan_uuid=?", column_name)
    //sql := "update `service_plans` set start_time=? where service_plan_uuid=?"
    fmt.Println(sql)
    _, rerr := engine.Exec(sql, start_time, service_plan_uuid)
    return rerr
}

func test_uuid() {
	engine, err := xorm.NewEngine("mysql", "root:restfulapi123@(10.209.224.161:10022)/wdf_backend_db")
	if err != nil {
		fmt.Println("error1")
	}
	defer engine.Close()
	engine.ShowSQL(true)
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)
	//modelbase := new(models.ModelBase)
	models := make([]models.ModelBase, 4)
	//models[0].Id =
	models[0].ModelUuid = u1.String()
	models[0].ModelType = "charg_type"
	models[0].TypeValue = "hour"
	models[1].ModelUuid = uuid.NewV4().String()
	models[1].ModelType = "charg_type"
	models[1].TypeValue = "minute"
	models[2].ModelUuid = uuid.NewV4().String()
	models[2].ModelType = "cluster_type"
	models[2].TypeValue = "kafka"
	models[3].ModelUuid = uuid.NewV4().String()
	models[3].ModelType = "cluster_type"
	models[3].TypeValue = "spark"
	affected, err := engine.Insert(&models)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("affected:", affected)
}

func test_join() {
	engine, err := xorm.NewEngine("mysql", "root:restfulapi123@(10.209.224.161:10022)/czq")
	if err != nil {
		fmt.Println("error1")
	}
	defer engine.Close()
	engine.ShowSQL(true)
	users := make([]ServiceModule, 0)
	engine.Join("INNER", "service_base", "module_base.service_id = service_base.service_id").Find(&users)
	for _, i := range users {
		fmt.Println("i:\n", i)
	}
}

type Services struct {
	Id                          int
	ServiceUuid                 string
	ServiceBaseUuid             string
	//UserId                      string
	ProjectUuid                 string
    Env                         string                 `json:"env"`
    CreateTime                  time.Time              `xorm:"created"`
	EndTime                     time.Time
	IfImmediate                 int
	IfShared                    int
	IfPersist                   int
	IfJobPlan                   int
	Status                      string
	ServiceAddr                 string
	AddrForUse                  string
	CpuAccumulate               int64
	MemoryAccumulate            int64
	StorageAccumulate           int64
	BillingAccumulate           float64
	Parameters                  string
}
type Modules struct {
	Id              int
	ModuleUuid      string
    ModuleName      string
	ServiceUuid     string
	ModuleBaseUuid  string
	ReplicasValue   int
	ReplicasUnit    string
	CpuValue        int
	CpuUnit         string
	MemoryValue     int
	MemoryUnit      string
	StorageValue    int
	StorageUnit     string
	ReadyReplicas   int
	CreatedReplicas int
}
type ModuleBase struct {
    Id                          int             `json:"id,omitempty"`
    ModuleBaseUuid              string          `json:"module_base_uuid,omitempty"`
    ModuleName                  string          `json:"module_name"`
    ModuleKind                  string          `json:"module_kind,omitempty"`
    Version                     string          `json:"version"`
    ServiceBaseUuid             string          `json:"service_base_uuid,omitempty"`
    ModuleParametersMapping     string          `json:"module_parameters_mapping,omitempty"`
    JobPlan                     int             `json:"job_plan,omitempty"`
}
type ModuleBaseWithServices struct {
    ModuleBase              `xorm:"extends"`
    Modules                 `xorm:"extends"`
    Services                `xorm:"extends"`
    Projects                `xorm:"extends"`
}
type Pvcs struct {
	Id               int
	PvcUuid          string
	PvcName          string
	ProjectName      string
	VolumeName       string
	VolumePath       string
	StorageClassName string
	ServiceUuid      string
	Storage          int
	StorageUnit      string
	StorageType      string
	CreateTime       time.Time
	RestoreTime      time.Time
	Status           string
	StorageUsage     float64
	Billing          float64
	RunTime          int
}
func (ModuleBaseWithServices) TableName() string {
    return "module_base"
}
func GetModuleBaseWithServices(service_uuid string) (*ModuleBaseWithServices, error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    mb := new(ModuleBaseWithServices)
    engine.Join("INNER", "modules", "modules.module_base_uuid = module_base.module_base_uuid").Join("INNER", "services", "module_base.service_base_uuid = services.service_base_uuid").Join("INNER", "projects", "services.project_uuid = projects.project_uuid").Where("services.service_uuid = ?", service_uuid).And("module_base.job_plan = ?", 1).Get(mb)
    return mb, nil
}
func GetSelector(mbs *ModuleBaseWithServices) ([]*string, error) {
    /*select_str := "oc get " + mbs.ModuleBase.ModuleKind
    select_str += " " + mbs.Modules.ModuleName
    select_str += " -n " + mbs.Projects.ProjectName
    select_str += " --template={{.spec.selector.matchLabels}}"*/
    //select_str := "oc get pod -n " + mbs.Projects.ProjectName
    //select_str += " -o name | grep " + mbs.Modules.ModuleName
    select_str := "oc get pod -n " + mbs.Projects.ProjectName
    select_str += " | grep " + mbs.Modules.ModuleName
    select_str += " | grep Running | awk '{print $1}'"
    fmt.Println(select_str)
    out, err := exec.Command("bash", "-c", select_str).CombinedOutput()
    fmt.Println(string(out))
    if err != nil {
        return nil, fmt.Errorf("RealResetService patch out:%s", string(out))
    }
    ret := make([]string, 0)
    first_arr := strings.Split(string(out), "\n")
    for _, tmp := range first_arr {
        fmt.Print(tmp)
        if len(tmp) != 0 {
            ret = append(ret, tmp)
        }
    }
    fmt.Println(ret)
    return nil, nil
}
type Test struct {
	Id   int
    A    bool        `xorm:"Bool"`
	Name string
}
func test_bool(table_name, uuid string) (error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)

    id_name := strings.TrimSuffix(table_name, "s") + "_uuid"
    sql := fmt.Sprintf("select `status` from `%s` where `%s`=?", table_name, id_name)
    res, err := engine.Query(sql, uuid)
    if err != nil || len(res) == 0 {
        return err
    }
    fmt.Println(string(res[0]["status"]))
    return nil
}
type InstancesWithPvcsByServiceId struct {
    Instances           `xorm:"extends"`
    Modules             `xorm:"extends"`
    Services            `xorm:"extends"`
    Pvcs                `xorm:"extends"`
}
type Instances struct {
	Id           int
	InstanceUuid string
	InstanceName string
	ProjectName  string
	ModuleUuid   string
	PvcUuid      string
	NodeIp       string
	PodIp        string
	StartTime    time.Time
	EndTime      time.Time
	Status       string
	CpuUsage     float64
	CpuLimit     int
	MemoryUsage  float64
	MemoryLimit  int
	RunTime      int
}
func (InstancesWithPvcsByServiceId) TableName() string {
    return "instances"
}
func test_limit(limit, offset int) (error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)


    sid := "91a62eb7-09bc-480f-a946-c08cf95fdb5a"
    sortby := "instance_name"
    var ret []*InstancesWithPvcsByServiceId
    _ = engine.Join("INNER", "modules", "instances.module_uuid = modules.module_uuid").Join("INNER", "services", "services.service_uuid = modules.service_uuid").Join("INNER", "pvcs", "instances.pvc_uuid  = pvcs.pvc_uuid").Where("services.service_uuid = ?", sid).Limit(limit, offset).Desc(strings.TrimPrefix(sortby, "-")).Find(&ret)

    count_ret := new(InstancesWithPvcsByServiceId)
    total, _ := engine.Join("INNER", "modules", "instances.module_uuid = modules.module_uuid").Join("INNER", "services", "services.service_uuid = modules.service_uuid").Join("INNER", "pvcs", "instances.pvc_uuid  = pvcs.pvc_uuid").Where("services.service_uuid = ?", sid).Count(count_ret)
    fmt.Println(len(ret))
    fmt.Println(total)
    return nil
}
func GetServicePlansUuids() (error) {
    dbaddr := "root:restfulapi123@(10.209.224.162:10456)/czq2?charset=utf8"
    engine, err := xorm.NewEngine("mysql", dbaddr)
    defer engine.Close()
	if err != nil {
        fmt.Printf("Connect db error:", err)
	}
    engine.SetMaxIdleConns(10)
	engine.ShowSQL(true)
    var service_uuids []string
    err = engine.Table("service_plans").Cols("service_plan_uuid").Where("service_plans.status = ?", "Created").Find(&service_uuids)
    if err != nil {
        return fmt.Errorf("GetProjectByUuid select db error:%s", err)
    }
    for _, id := range service_uuids {
        fmt.Printf("id:%s\n", id)
    }
    return nil
}
func test_time_dx(year, month, day, hour, minute int) bool {
    str_month := strconv.Itoa(month)
    if month < 10 {
        str_month = "0" + str_month
    }
    str_day := strconv.Itoa(day)
    if day < 10 {
        str_day = "0" + str_day
    }
    str_hour := strconv.Itoa(hour)
    if hour < 10 {
        str_hour = "0" + str_hour
    }
    str_minute := strconv.Itoa(minute)
    if minute < 10 {
        str_minute = "0" + str_minute
    }
    first_run_time := fmt.Sprintf("%d-%s-%s %s:%s:00", year, str_month, str_day, str_hour, str_minute)
    /*timeLayout := "2006-01-02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    firsRunTime, err := time.ParseInLocation(timeLayout, first_run_time, loc)
    if err != nil {
        fmt.Println("parse in location error\n")
        return false
    }*/
    now_str := time.Now().Format("2006-01-02 15:04:05")
    fmt.Println(now_str)
    fmt.Println(first_run_time)
    if now_str < first_run_time {
        fmt.Println("first run time is less than now")
    } else {
        fmt.Println("ok")
    }
    return true

}
func main() {
	//test_user()
	//test_join()
//	test_uuid()
    //UpdateServicePlanTable("9fccef38-a315-4d75-b745-0453e83a6f23", "start_time", "12:00")
    /*mb, _ := GetModuleBaseWithServices("91a62eb7-09bc-480f-a946-c08cf95fdb5a")
    GetSelector(mb)*/
    /*c := "oc rsh accp-d-1744562428-2p3xm ps -efw | grep test.sh"
    _, err := exec.Command("bash", "-c", c).CombinedOutput()
    if err != nil {
        fmt.Println("1212121")
    } else {
        fmt.Println("333333")
    }
    i := 5
    a := -1
    i += a
    fmt.Print(i)
    depend_service_arr := strings.Split("", ",")
    for _, depend_service_uuid := range depend_service_arr {
        fmt.Printf("aa:", depend_service_uuid)
    }*/
/*    str := "sample.xml"
	cmname := strings.Replace(str, ".", "\\.", -1)
    fmt.Printf(cmname)
    test_bool("services", "157241fa-2b12-4a25-90df-4cd90dd699bb")*/
//    test_limit(5, 0)
    /*err := CommunicateProjects("f18214e1-7b95-40d3-9596-54a49b33bfba", "zk0-1.zk.czq-project.svc", true)
    if err != nil {
        fmt.Println(err.Error())
    }*/
//    GetServicePlansUuids()
    test_time_dx(2017, 12, 11, 15, 56)
}
