package good

import "github.com/vaberof/hezzl-backend/pkg/domain"

type GoodStorage interface {
	Create(projectId domain.ProjectId, name domain.GoodName) (*Good, error)
	Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*Good, error)
	Delete(id domain.GoodId, projectId domain.ProjectId) (*Good, error)
	List(limit, offset int) ([]*Good, error)
	ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority domain.GoodPriority) ([]*Good, error)
	IsExists(id domain.GoodId, projectId domain.ProjectId) (bool, error)
}
