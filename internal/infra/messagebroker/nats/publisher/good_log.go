package publisher

import (
	"time"
)

type GoodLog struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"eventTime"`
}
