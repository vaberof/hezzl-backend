package good

import (
	"github.com/vaberof/hezzl-backend/pkg/domain"
)

type Good struct {
	Id          domain.GoodId
	ProjectId   domain.ProjectId
	Name        domain.GoodName
	Description domain.GoodDescription
	Priority    domain.GoodPriority
	Removed     domain.GoodRemoved
	CreatedAt   domain.GoodCreatedAt
}
