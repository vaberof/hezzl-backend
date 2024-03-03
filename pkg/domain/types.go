package domain

import "time"

type GoodId int64

type GoodName string

func (name *GoodName) String() string {
	return string(*name)
}

type GoodDescription string

func (description *GoodDescription) String() string {
	return string(*description)
}

type GoodPriority int

type GoodRemoved bool

type GoodCreatedAt time.Time

type ProjectId int64
