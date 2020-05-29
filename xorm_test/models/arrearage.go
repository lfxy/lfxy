package models

import (
	"time"
)

type Arrearage struct {
	Id            int
	ArrearageUuid string
	UserId        string
	CreateTime    time.Time
	RecoverTime   time.Time
}
