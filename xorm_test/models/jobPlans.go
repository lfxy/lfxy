package models

import (
	"time"
)

type JobPlans struct {
	Id                int
	JobPlanUuid       string
	JobPlanName       string
	ServiceUuid       string
	CreateTime        time.Time
	DeleteTime        time.Time
	PeriodType        string
	StartTime         time.Time
	EndTime           time.Time
	Status            string
	Day               int
	Hour              int
	Minute            int
	CronExpression    string
	ShellScripts      string
	DependServiceUuid string
}
