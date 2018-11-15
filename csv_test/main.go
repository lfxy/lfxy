package main

import (
		"encoding/csv"
        "fmt"
		"io/ioutil"
		"strings"
        "encoding/json"
)

type ttype struct {
	v1	int
	v2	float32
}

func test(){
        fileName := "c3.csv"
        cntb,err := ioutil.ReadFile(fileName)
        if err != nil {
			fmt.Println(err)
        }
        r2 := csv.NewReader(strings.NewReader(string(cntb)))
        ss,_ := r2.ReadAll()
        fmt.Println(ss)
        sz := len(ss)
		fmt.Println("----------sz:", sz)
		num := -1
		name := ""
		for num, name = range ss[0] {
			if name == "req_tot" {
				break
			}
		}
        fmt.Println("----------num", num)
        for i:=1;i<sz;i++{
			fmt.Println(ss[i][num])
        }

		fmt.Println("=============================================")
		cm := make(map[string]ttype)
		tt := ttype {
			v1:		10,
			v2:		20,
		}
		cm["aadd"] = tt

		fmt.Println(cm["aadd"].v2)

}

func teststr(){
    t := [3]string{"aaa", "bbb", "ccc"}
    var s string
    for _, i := range t {
        s += "name != " + i + " and "
    }
    fmt.Println(s)
    s = strings.TrimSuffix(s, " and ")
    fmt.Println(s)

}


type ModuleConfig struct {
    Name            string          `json:"name"`
    Value           int             `json:"value"`
    ValueRange      string          `json:"value_range"`
    ValueUnit       string          `json:"value_unit"`
    DisplayBox      string          `json:"display_box"`
    Describe        string          `json:"describe"`
}

type mtest struct {
    Modules         map[string][]*ModuleConfig            `json:"modules"`
}

func main(){
    strt := `{"modules":{"spark-master":[{"name":"replicas","value":1,"value_range":"1,       10","value_unit":"number","display_box":"dropdown","describe":"replicas"},{"name":"cpu","value":1,"value_range":"1,       16","value_unit":"core","display_box":"dropdown","describe":"cpu"}]}}`
    var servicesfDeploy mtest
    err := json.Unmarshal([]byte(strt), &servicesfDeploy)
    if err != nil {
        fmt.Println("err:", err)
        return
    }
    fmt.Println("content:", servicesfDeploy)
    for k, v := range servicesfDeploy.Modules {
        fmt.Println("k:", k)
        for _, c := range v {
            fmt.Println("name:", c.Name)
        }

    }
}
