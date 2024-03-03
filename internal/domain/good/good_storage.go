package good

import "github.com/vaberof/hezzl-backend/pkg/domain"

type GoodStorage interface {
	Create(projectId domain.ProjectId, name domain.GoodName) (*Good, error)
	Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*Good, error)
	Delete(id domain.GoodId, projectId domain.ProjectId) error
	List(limit, offset int) ([]*Good, error)
	ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority int) ([]*Good, error)
	IsExists(id domain.GoodId, projectId domain.ProjectId) (bool, error)
}
