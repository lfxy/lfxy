package main

import (
//	"bytes"
    "fmt"
    "github.com/influxdata/influxdb1-client/v2"
	"log"
	"os"
//    "os/exec"
//  "strconv"
//  "strings"
  "time"
)

const (
	Mydb = "bigdata-db"
	username = "bigdatawrite"
	password = "agoradevops2020"
	url = "http://bigdata.influx.agoralab.co/"
)

type DbClient struct {
	Cli client.Client
}

/*var influxCli DbClient

func init(){
	influxCli = DbClient{ DbConfig() }
}*/

func DbConfig()client.Client{
	Cli,err := client.NewHTTPClient(client.HTTPConfig{Addr:url, Username: username, Password:password,})
	if err != nil {
		printlog(fmt.Sprintf("%s",err))
	}
	return Cli
}

func (c DbClient)WritePoint(measurement, tag string, fields map[string]interface{}){
	bp,err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: Mydb,
		Precision: "s",
        RetentionPolicy: "rp_30d",
	})
	if err != nil {
		printlog(fmt.Sprintf("%s",err))
	}
	/*lag,err := GetLag(groupid)
	if err != nil  {
		printlog(fmt.Sprintf("%s",err))
		return
	}*/
	tags := map[string]string{"application": tag}
	pt,err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		printlog(fmt.Sprintf("%s",err))
	}
	bp.AddPoint(pt)
	err = c.Cli.Write(bp)
	if err != nil {
		printlog(fmt.Sprintf("%s",err))
	}else {
		printlog("insert success")
	}

}

func printlog(a string){
	f, err := os.OpenFile("lag_monitor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(a)
}


/*func main(){
    var influxCli *DbClient = new(DbClient)
    influxCli.Cli = DbConfig()
	influxCli.WritePoint("yarn_crontab_loadaverage", "attributor", 11)

    */

/*	t := time.NewTicker(time.Second*10)
	defer t.Stop()
	for{
		select {
		case <-t.C:
			influxCli.WritePoint("attributor")
			influxCli.WritePoint("MessageHistoryNADiv0")
			if Mydb == "rtm_kafka_cn"{
				influxCli.WritePoint("MessageHistoryCNDiv1")
			}
		}
	}*/
//}

