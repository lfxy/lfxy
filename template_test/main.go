package main

import (
	"fmt"
	"strings"
    "text/template"
    "log"
    "os"
    "os/exec"
    "encoding/json"
    "strconv"
    "time"
    "math"
)

const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.
{{- else}}
It is a shame you couldn't make it to the wedding.
{{- end}}
{{with .Gift -}}
Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`
const cmdtext = `
#!/bin/bash

cat {{.File1}}
cat {{.File2}}
`

// Prepare some data to insert into the template.
type Recipient struct {
    Name, Gift string
    Attended   bool
}
var recipients = []Recipient{
    {"Aunt Mildred", "bone china tea set", true},
    {"Uncle John", "moleskin pants", false},
    {"Cousin Rodney", "", false},
}

// Create a new template and parse the letter into it.
func test_template() {
	t := template.Must(template.ParseFiles("./test.txt"))

	// Execute the template for each recipient.
	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}

}

type ArgsOfScript struct {
    Args    map[string]string
}

type Testarg struct {
    File1   string
    File2   string
}

func test_exec() {
    f, err := os.Create("./cmd_tmp.sh")
    defer f.Close()
    defer os.Remove("./cmd_tmp.sh")

    /*para1 := "File1"
    para2 := "File2"
    v1 := "a.txt"
    v2 := "b.txt"
    jstr := "{"
    jstr += "\"" + para1 + "\":\"" + v1 + "\","
    jstr += "\"" + para2 + "\":\"" + v2 + "\""
    jstr += "}"
    fmt.Println("jstr:", jstr)*/
    var strargs = []byte(`{
        "File1":"a.txt",
        "File2":"b.txt",
    }`)
    //var strargs = []byte(jstr)

    var animals interface{}
    err = json.Unmarshal(strargs, &animals)
    if err != nil {
        fmt.Println("json error : ", err)
    }
    t := template.Must(template.ParseFiles("./cmd.sh"))
    fmt.Printf("animals:%+v", animals)
    ta := Testarg{"a.txt","b.txt"}
    fmt.Printf("ta:%+v", ta)
    err = t.Execute(f, animals)
    if err != nil {
        log.Println("executing template:", err)
    }
    f.Chmod(0777)


	out, err := exec.Command("bash", "-c", "./cmd_tmp.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The date is %s\n", out)

}

func test_exec_string() {
    var strargs1 = []byte(`{
        "File1":"a.txt"
    }`)
    var strargs2 = []byte(`{
        "File2":"b.txt"
    }`)

    f, err := os.Create("./cmd_tmp.sh")
    defer f.Close()
    defer os.Remove("./cmd_tmp.sh")

    var animals1 interface{}
    var animals2 interface{}
    err = json.Unmarshal(strargs1, &animals1)
    err = json.Unmarshal(strargs2, &animals2)
    if err != nil {
        fmt.Println("json error : ", err)
    }
    fmt.Printf("animals1:%+v", animals1)
    fmt.Printf("animals2:%+v", animals2)
	t := template.Must(template.New("cmdtext").Parse(cmdtext))
    err = t.Execute(f, animals1)
    //t1 := template.Must(template.ParseFiles("./cmd_tmp.sh"))
    //err = t1.Execute(f, animals2)
    if err != nil {
        log.Println("executing template:", err)
    }
    f.Chmod(0777)


	out, err := exec.Command("bash", "-c", "./cmd_tmp.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The date is %s\n", out)

}

func test_json_combine(){
    var json1 = [...]string{`{"aaa":"111","bbb":"222","ccc":"333"}`, `{"ddd":"444","eee":"555","fff":"666"}`, `{"ggg":"777","hhh":"888","iii":"999"}`}
	//fmt.Printf(json1[0], json1[0])
    jsonstr := "{"
    for _, obj := range json1 {
        jsonstr += strings.TrimSuffix(strings.TrimPrefix(obj, "{"), "}")
        jsonstr += ","
    }
    jsonstr = strings.TrimSuffix(jsonstr, ",")
    jsonstr += "}"
    fmt.Printf(jsonstr)
}

func test_json_parse() {
    json1 := []byte(`{"Aaa":"111","Bbb":"222","Ccc":"333"}`)
    var animals1 interface{}
    json.Unmarshal(json1, &animals1)
    if ans, ok := animals1.(map[string]interface{}); ok {
        for k, v := range ans {
            if b, ok := v.(string); ok {
                fmt.Println(k)
                fmt.Println(b)
            }
        }
    } else {
        fmt.Println("error")
    }
}

func test_var(src string, args ...interface{}) {
    fmt.Println(src)
    fmt.Println(args[0])
    fmt.Println(args[1])
}

func test_stringtoint(){
    memory := "4"
    memory_int, err := strconv.Atoi(memory)
    if err != nil {
        fmt.Println("error")
    }
    xmx_int := float64(memory_int) * 0.7
    fmt.Println(xmx_int)
    xmx_str := strconv.Itoa(int(xmx_int))
    fmt.Println(xmx_str)
}

func Round(f float64, n int) float64 {
    pow10_n := math.Pow10(n)
    return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
func test_time(){
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    fmt.Println(time.Now().String())
    fa := 1.44654
    fmt.Println(Round(fa, 2))
	out, err := exec.Command("bash", "-c", "ls -l").Output()
	if err != nil {
        fmt.Printf("err %s\n", err)
	}
	fmt.Printf("The date is %s\n", out)
}
func test_exec_err(){
	out, err := exec.Command("bash", "-c", "oc get statefulsets spark-worker-mj2 -n czq-project").Output()
    if err != nil {
        fmt.Print("11111111:", err)
    }
    fmt.Print(string(out))
	out, err = exec.Command("bash", "-c", "oc get statefulsets spark-worker-mj23 -n czq-project").CombinedOutput()
    if err != nil {
        fmt.Print("22222222:", err)
    }
    fmt.Print(string(out))

}


func GetConfigMapName(module_kind, module_name, project_name string) (string, error) {
	checkcmd := "oc get " + module_kind
	checkcmd += " " + module_name
	checkcmd += " -n " + project_name
    checkcmd += " -o jsonpath=\"{.spec.template.spec.volumes[*].configMap.name}\""
	fmt.Printf("%s\n", checkcmd)
	name_out, err := exec.Command("bash", "-c", checkcmd).CombinedOutput()
    fmt.Printf("%s\n", string(name_out))
    if err != nil {
        fmt.Println("GetConfigMaps get configMap name error:%s", err.Error())
        return "", err
    }
    return string(name_out), nil
}

func GetConfigMaps(module_kind, module_name, project_name string) ([]map[string]string, error) {
    /*m_ret := make(map[string]string)
	basecmd := "oc get " + module_kind
	basecmd += " " + module_name
	basecmd += " -n " + project_name
    namecmd := basecmd + " -o jsonpath=\"{.spec.template.spec.volumes[*].configMap.name}\""
	fmt.Printf("%s\n", namecmd)
	name_out, err := exec.Command("bash", "-c", namecmd).CombinedOutput()
    //fmt.Printf("%s\n", string(name_out))
    if err != nil {
        fmt.Println("GetConfigMaps get configMap name error:%s", err.Error())
        return nil, err
    }
    m_ret["cm_name"] = string(name_out)

    keycmd := basecmd + " -o jsonpath=\"{.spec.template.spec.volumes[*].configMap.items[*].key}\""
	fmt.Printf("%s\n", keycmd)
	key_out, err := exec.Command("bash", "-c", keycmd).CombinedOutput()
    if err != nil {
        fmt.Println("GetConfigMaps get configMap key error:%s", err.Error())
        return nil, err
    }
    m_ret["cm_key"] = string(key_out)

    pathcmd := basecmd + " -o jsonpath=\"{.spec.template.spec.volumes[*].configMap.items[*].path}\""
	fmt.Printf("%s\n", pathcmd)
    path_out, err := exec.Command("bash", "-c", pathcmd).CombinedOutput()
    if err != nil {
        fmt.Println("GetConfigMaps get configMap path error:%s", err.Error())
        return nil, err
    }
    m_ret["cm_path"] = string(path_out)


	valuename := strings.Replace(m_ret["cm_key"], ".", "\\.", -1)
    valuecmd := "oc get configmaps " + m_ret["cm_name"]
	valuecmd += " -n " + project_name
    valuecmd += " -o jsonpath=\"{.data."
    valuecmd += valuename + "}\""
	fmt.Printf("%s\n", valuecmd)
	value_out, err := exec.Command("bash", "-c", valuecmd).CombinedOutput()
    if err != nil {
        fmt.Println("GetConfigMaps get configMap value error:%s", err.Error())
        return nil, err
    }
    m_ret["cm_value"] = string(value_out)
    return m_ret, nil
    /*tmp_out := strings.TrimPrefix(string(cm_out), "map[")
    tmp_out = strings.TrimSuffix(tmp_out, "]")
    tmp_arr := strings.SplitN(tmp_out, ":", 2)
    fmt.Printf("key:%s\nvalue:%s\n", tmp_arr[0], tmp_arr[1])
    return tmp_arr[0], tmp_arr[1], nil*/
    //m_ret := make(map[string]string)
	basecmd := "oc get " + module_kind
	basecmd += " " + module_name
	basecmd += " -n " + project_name
    namecmd := basecmd + " -o jsonpath=\"{range .spec.template.spec.volumes[*]}{.configMap.name}\t{end}\""
	fmt.Printf("%s\n", namecmd)
	name_out, err := exec.Command("bash", "-c", namecmd).CombinedOutput()
    //fmt.Printf("%s\n", string(name_out))
    if err != nil || string(name_out) == "" {
        fmt.Println("GetConfigMaps get configMap name error:%s", err.Error())
        return nil, err
    }

    name_arr := strings.Split(string(name_out), "\t")
    count := 0
    for _, names := range name_arr {
        if names != "" {
            count++
        }
    }
    fmt.Printf("count111:%d\n", count)
    m_arr := make([]map[string]string, count)
    count = 0
    for _, names := range name_arr {
        if names != "" {
            m_tmp := make(map[string]string)
            m_tmp["cm_name"] = names
            m_arr[count] = m_tmp
            fmt.Printf("count 222:%d\n", count)
            count++
        }
    }
    //m_ret["cm_name"] = string(name_out)

    keycmd := basecmd + " -o jsonpath=\"{range .spec.template.spec.volumes[*]}{.configMap.items[*].key}\t{end}\""
	fmt.Printf("%s\n", keycmd)
	key_out, err := exec.Command("bash", "-c", keycmd).CombinedOutput()
    if err != nil || string(key_out) == "" {
        fmt.Println("GetConfigMaps get configMap key error:%s", err.Error())
        return nil, err
    }
    key_arr := strings.Split(string(key_out), "\t")
    fmt.Println(key_arr)
    count = 0
    for _, keys := range key_arr {
        if keys != "" {
            fmt.Printf("count 333:%d, keys:%s\n", count, keys)
            m_arr[count]["cm_key"] = keys
            count++
        }
    }
    //m_ret["cm_key"] = string(key_out)

    pathcmd := basecmd + " -o jsonpath=\"{range .spec.template.spec.volumes[*]}{.configMap.items[*].path}\t{end}\""
	fmt.Printf("%s\n", pathcmd)
    path_out, err := exec.Command("bash", "-c", pathcmd).CombinedOutput()
    if err != nil || string(path_out) == "" {
        fmt.Println("GetConfigMaps get configMap path error:%s", err.Error())
        return nil, err
    }
    path_arr := strings.Split(string(path_out), "\t")
    count = 0
    for _, paths := range path_arr {
        if paths != "" {
            fmt.Printf("count 444:%d\n", count)
            m_arr[count]["cm_path"] = paths
            count++
        }
    }
    //m_ret["cm_path"] = string(path_out)


    for i, m_ret := range m_arr {
        valuename := strings.Replace(m_ret["cm_key"], ".", "\\.", -1)
        valuecmd := "oc get configmaps " + m_ret["cm_name"]
        valuecmd += " -n " + project_name
        valuecmd += " -o jsonpath=\"{.data."
        valuecmd += valuename + "}\""
        fmt.Printf("%s\n", valuecmd)
        value_out, err := exec.Command("bash", "-c", valuecmd).CombinedOutput()
        if err != nil || string(value_out) == "" {
            fmt.Println("GetConfigMaps get configMap value error:%s", err.Error())
            return nil, err
        }
        m_ret["cm_value"] = string(value_out)
        m_arr[i] = m_ret
    }
    return m_arr, nil
}
func GetEnvs(module_kind, module_name, project_name string) (map[string]string, error) {
    m_ret := make(map[string]string)
	basecmd := "oc get " + module_kind
	basecmd += " " + module_name
	basecmd += " -n " + project_name
    envcmd := basecmd + " -o jsonpath=\"{range .spec.template.spec.containers[0].env[*]}{.name}\t{.value}\n{end}\""
	fmt.Printf("%s\n", envcmd)
    env_out, err := exec.Command("bash", "-c", envcmd).CombinedOutput()
    if err != nil {
        fmt.Println("GetConfigMaps get configMap env error:%s", err.Error())
        return nil, err
    }
    if string(env_out) == "" || strings.Contains(string(env_out), "not found") {
        return nil, fmt.Errorf("GetEnvs err")
    }
    env_tmp_arr := strings.Split(string(env_out), "\n")
    for _, env_tmp := range env_tmp_arr {
        //fmt.Printf("i:%d, s:%s\n", i, env_tmp)
        if env_tmp != "" {
            envs := strings.Split(env_tmp, "\t")
            /*if len(envs) != 2 {
                return nil, fmt.Errorf("GetEnvs len error")
            }*/
            m_ret[envs[0]] = envs[1]
        }
    }
    return m_ret, nil
}
type EnvParam struct {
    Name        string          `json:"name"`
    Value       string          `json:"value"`
}
func test_envs(module_kind, module_name, project_name string) (error) {
    m_ret := make(map[string]string)
    m_ret["JVM_OPTS"] = "-Xmx8000m -Xms8000m"
    m_ret["OPTION_LIBS"] = "ignite-kubernetes"
    m_ret["CONFIG_URI"] = "file:////opt/ignite/apache-ignite-fabric-2.2.0-bin/config/example-kube.xml"
    ret_i, err := json.Marshal(m_ret)
    if err != nil  {
        fmt.Println("error!")
    }
    fmt.Printf("ret:%s\n", string(ret_i))

    env_arr := make([]*EnvParam, 0)
    for k, v := range m_ret {
        envobj := new(EnvParam)
        envobj.Name = k
        envobj.Value = v
        env_arr = append(env_arr, envobj)
    }
    ret3, err := json.Marshal(env_arr)
    if err != nil  {
        fmt.Println("error3!")
    }
    fmt.Printf("ret3:%s\n", string(ret3))

    ret2 := "["
    for k, v := range m_ret {
        ret2 += "{"
        ret2 += "\"name\":\""
        ret2 += k
        ret2 += "\",\"value\":\""
        ret2 += v
        ret2 += "\"},"
    }
    ret2 = strings.TrimSuffix(ret2, ",")
    ret2 += "]"
    fmt.Printf("ret2:%s\n", string(ret2))
    namecmd := "oc get " + module_kind
    namecmd += " " + module_name
    namecmd += " -n " + project_name
	namecmd += " -o jsonpath=\"{.spec.template.spec.containers[*].name}\""
    fmt.Printf("namecmd:%s\n", namecmd)
    container_name, err := exec.Command("bash", "-c", namecmd).CombinedOutput()
    if err != nil {
        return err
    }
    envpathstr := fmt.Sprintf("{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"env\":%s}]}}}}", string(container_name), ret3)
    envcmd := "oc patch " + module_kind
    envcmd += " " + module_name
    envcmd += " -n " + project_name
    envcmd += " -p '" + envpathstr
    envcmd += "'"
    fmt.Printf("namecmd:%s\n", envcmd)
    patch_out, err := exec.Command("bash", "-c", envcmd).CombinedOutput()
    fmt.Printf("patch out:%s", string(patch_out))
    if err != nil {
        return err
    }
    return nil
}
func test_service_add() {
    scriptname := "./test_addr.sh"
	out, err := exec.Command("bash", "-c", scriptname).CombinedOutput()
    fmt.Printf("out:%s\n", string(out))
    if err != nil {
        fmt.Printf("error:%s\n", err.Error())
        return
    }
    strout := string(out)
    if strings.Contains(strout, "service_addr") {
        out_arr := strings.Split(strout, "\n")
        for _, tmp_addr := range out_arr {
            if strings.Contains(tmp_addr, "service_addr") {
                //tmp_addr = "service_addr\taaa"
                addr_str_arr := strings.Split(tmp_addr, " ")
                fmt.Printf("tmp_addr===:%s\n", tmp_addr)
                fmt.Printf("addr_str_arr===:%s\n", addr_str_arr)
                if len(addr_str_arr) != 2{
                    fmt.Printf("RealDeploy error parse service_addr:%s", tmp_addr)
                    break
                }
                fmt.Printf("===:%s", addr_str_arr[1])
            }
        }
    }
}
func main(){
    //test_template()
    //test_json_combine()
    //test_json_parse()
//    test_var("ssssrc", "aa", "bbb")
    //test_exec()
    //test_stringtoint()
//    test_time()
    //test_exec_err()
	/*out, err := exec.Command("bash", "-c", "oc whoami -t").CombinedOutput()
    if err != nil {

    }
    fmt.Printf("%v", string(out))*/

    /*m_ret, err := GetConfigMaps("statefulsets", "ignite-9j7", "czq-project")
    if err != nil {
        return
    }
    for i, m_ret := range m_ret {
        fmt.Printf("i:%d\n", i)
        for k, v := range m_ret {
            fmt.Printf("key:%s\nvalue:%s\n", k, v)
        }
    }*/

    /*m_ret, err := GetEnvs("statefulsets", "ignite-9j7", "czq-project")
    if err != nil{
        fmt.Print(err.Error())
        return
    }
    for k, v := range m_ret {
        fmt.Printf("key:%s\nvalue:%s\n", k, v)
    }*/
    //test_envs("statefulsets", "ignite-9j7", "czq-project")
    test_service_add()
}
