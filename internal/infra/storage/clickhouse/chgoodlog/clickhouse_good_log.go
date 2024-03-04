package chgoodlog

import "time"

type GoodLog struct {
	Id          int64
	ProjectId   int64
	Name        string
	Description string
	Priority    int
	Removed     bool
	EventTime   time.Time
}
