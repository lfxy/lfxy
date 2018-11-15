package main

import (
        "fmt"
		"time"
)

type MetricType int8
type ValueType int8
type MetricValue struct {
	IntValue    int64
	FloatValue  float32
	CustomValue float64
	MetricType  MetricType
	ValueType   ValueType
}
type MetricSet struct {
	CreateTime     time.Time
	ScrapeTime     time.Time
	MetricValues   map[string]MetricValue
	Labels         map[string]string
}

type DataBatch struct {
	Timestamp time.Time
	MetricSets map[string]*MetricSet
}
func ScrapeMetrics() *DataBatch {
	fmt.Println("ScrapeMetrics 00000")
	result := &DataBatch{
		Timestamp:  time.Now(),
		MetricSets: map[string]*MetricSet{},
	}
	if result.MetricSets == nil {
		fmt.Println("nil 1111111111")
	}

	return result
}

func main(){
	customResponseChannel := make(chan *DataBatch)
	tt := 20 * time.Second
	startTime := time.Now()
	timeoutTime := startTime.Add(tt)
	go func(channel chan *DataBatch, timeoutTime time.Time) {
		customemetrics := ScrapeMetrics()
		now := time.Now()
		if customemetrics == nil {
			fmt.Println("nil 1111111111")
			//return
		}
		if len(customemetrics.MetricSets) == 0 {
			fmt.Println("nil 22222222")
			//return
		}

		timeForResponse := timeoutTime.Sub(now)
		select {
		case channel <- customemetrics:
			fmt.Println("success 1111111")
			// passed the response correctly.
			return
		case <-time.After(timeForResponse):
			fmt.Println("nil 33333333333")
			return
		}
	}(customResponseChannel, timeoutTime)

	now := time.Now()
	select {
	case customDataBatch := <-customResponseChannel:
		if customDataBatch != nil && len(customDataBatch.MetricSets) > 0 {
			for _, h_value := range customDataBatch.MetricSets {
				for _, l_v := range h_value.MetricValues {
					fmt.Println("czq3 realScrapeMetrics l_v:%s ", l_v.FloatValue)
				}
			}
		} else {
			fmt.Println("2222  =====0     ")
		}

	case <-time.After(timeoutTime.Sub(now)):
		fmt.Printf("timeout 22222222")
	}

}
