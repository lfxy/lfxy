package main

import (
    "fmt"
    //"os/exec"
    "flag"
    //"strings"
    //"strconv"
	//"sort"
    "time"
    "github.com/tsuna/gohbase"
    //"influxdbclient"
    "github.com/parnurzeal/gorequest"
    //"os"
    //"log"
    //"path/filepath"
    "encoding/json"
    //"encoding/csv"
    //"sync"
    //"runtime"
)
type RegionServerApp struct {
    Beans            RegionServerInfos             `json:"beans"`
}

type RegionServerInfos []RegionServerServer

type RegionServerServer struct {
    Name                             string               `json:"name"`
    Mutate99ThPercentile             int               `json:"Mutate_99th_percentile"`
}

func GetRegionLatency(zkips string, threashold int, latencywithtimes map[string]int) (int, error) {
    high_latency_num := 0
    client := gohbase.NewAdminClient(zkips)
    clusterstatus, err := client.ClusterStatus()
    if err != nil {
        return -1, fmt.Errorf("get cluster status err :%v", err)
    }
    //fmt.Printf("master: %s, %d\n", clusterstatus.GetMaster().GetHostName(), clusterstatus.GetMaster().GetPort())
    liveserves := clusterstatus.GetLiveServers()
	request := gorequest.New()
    for _, server := range liveserves {
        jmx_url := fmt.Sprintf("http://%s:60030/jmx?qry=Hadoop:service=HBase,name=RegionServer,sub=Server", server.GetServer().GetHostName())
        _, body, errs := request.Get(jmx_url).End()
        if len(errs) != 0 {
            return -1, fmt.Errorf("http get error:%v", errs)
        }
        //fmt.Println(body)
        regionserverapp := new(RegionServerApp)
        err = json.Unmarshal([]byte(body), regionserverapp)
        if err != nil {
            return -1, fmt.Errorf("parse json error :%v", err)
        }
        //fmt.Printf("%s, %d\n", server.GetServer().GetHostName(), server.GetServer().GetPort())
        if regionserverapp.Beans[0].Mutate99ThPercentile > threashold {
            //fmt.Printf("GetRegionLatency   %s:%d\n", server.GetServer().GetHostName(), regionserverapp.Beans[0].Mutate99ThPercentile)
            high_latency_num += 1

            latencytimes, s_exists := latencywithtimes[server.GetServer().GetHostName()]
            if s_exists {
                latencywithtimes[server.GetServer().GetHostName()] = latencytimes + 1
            } else {
                latencywithtimes[server.GetServer().GetHostName()] = 1
            }
        } else {
            latencywithtimes[server.GetServer().GetHostName()] = 0
        }
    }

    return high_latency_num, nil
}
func main(){
    zkips := flag.String("zks", "10.10.4.41,10.10.4.42,10.10.7.213,10.10.7.214,10.10.7.215", "zk ip address")
    tag := flag.String("tag", "newcm.hbase", "tag name")
    threash_value := flag.Int("value", 10, "load threashold")
    flag.Parse()
    //fmt.Printf("%s, %d\n", *zkips, *threash_value)
    latencywithtimes := make(map[string]int, 0)
	t := time.NewTicker(time.Second * 30)
	defer t.Stop()
	for{
		select {
		case <-t.C:
            //fmt.Println("...........")

            ret_num, err := GetRegionLatency(*zkips, *threash_value, latencywithtimes)
            if err != nil {
                fmt.Printf("main error :%v", err)
                continue
            }

            threetimelatencys := 0
            for k, v := range latencywithtimes {
                if v > 3 {
                    fmt.Printf("%s:%d\n", k, v)
                    threetimelatencys++
                }
            }

            var influxCli *DbClient = new(DbClient)
            influxCli.Cli = DbConfig()
            fields := map[string]interface{}{
                "latency_nodes":   ret_num,
                "three_times_latency_nodes":   threetimelatencys,
            }
            influxCli.WritePoint("hbase_crontab_latency", *tag, fields)
		}
	}
}
