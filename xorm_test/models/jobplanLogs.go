package models

import ()

type JobplanLogs struct {
	Id          int
	PlanLogUuid string
	JobPlanUuid string
	Logs        string
}
