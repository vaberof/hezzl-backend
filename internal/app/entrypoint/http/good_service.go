package http

import (
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/pkg/domain"
)

type GoodService interface {
	Create(projectId domain.ProjectId, name domain.GoodName) (*good.Good, error)
	Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*good.Good, error)
	Delete(id domain.GoodId, projectId domain.ProjectId) (*good.Good, error)
	List(limit, offset int) ([]*good.Good, error)
	ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority domain.GoodPriority) ([]*good.Good, error)
}
