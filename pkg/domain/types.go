package domain

import "time"

type GoodId int

func (goodId *GoodId) Int64() int64 {
	return int64(*goodId)
}

type GoodName string

func (name *GoodName) String() string {
	return string(*name)
}

type GoodDescription string

func (description *GoodDescription) String() string {
	if description == nil {
		return ""
	}
	return string(*description)
}

type GoodPriority int

func (goodPriority *GoodPriority) Int() int {
	return int(*goodPriority)
}

type GoodRemoved bool

func (goodRemoved *GoodRemoved) Bool() bool {
	return bool(*goodRemoved)
}

type GoodCreatedAt time.Time

func (goodCreatedAt *GoodCreatedAt) Time() time.Time {
	return time.Time(*goodCreatedAt)
}

type ProjectId int64

func (projectId *ProjectId) Int64() int64 {
	return int64(*projectId)
}
