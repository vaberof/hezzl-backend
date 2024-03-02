package pggood

import "time"

type Good struct {
	Id          int64
	ProjectId   int64
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}
