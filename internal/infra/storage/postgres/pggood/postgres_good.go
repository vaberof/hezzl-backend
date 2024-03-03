package pggood

import (
	"database/sql"
	"time"
)

type Good struct {
	Id          int64
	ProjectId   int64
	Name        string
	Description sql.NullString
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}
