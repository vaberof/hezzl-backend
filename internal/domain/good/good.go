package good

import (
	"github.com/vaberof/hezzl-backend/pkg/domain"
	"time"
)

type Good struct {
	Id          domain.GoodsId
	ProjectId   domain.ProjectId
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}
