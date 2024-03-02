package good

import "github.com/vaberof/hezzl-backend/pkg/domain"

type GoodStorage interface {
	Create(projectId domain.ProjectId, name string) (*Good, error)
	Update(id domain.GoodsId, projectId domain.ProjectId, name, description string) (*Good, error)
	Delete(id domain.GoodsId, projectId domain.ProjectId) error
	List(limit, offset int) ([]*Good, error)
	ChangePriority(id domain.GoodsId, projectId domain.ProjectId) ([]*Good, error)
	IsExists(id domain.GoodsId, projectId domain.ProjectId) (bool, error)
}
