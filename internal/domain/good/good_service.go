package good

import (
	"errors"
	"github.com/vaberof/hezzl-backend/pkg/domain"
)

var (
	ErrGoodNotFound = errors.New("")
)

type GoodService interface {
	Create(projectId domain.ProjectId, name string) (*Good, error)
	Update(id domain.GoodsId, projectId domain.ProjectId, name, description string) (*Good, error)
	Delete(id domain.GoodsId, projectId domain.ProjectId) error
	List(limit, offset int) ([]*Good, error)
	ChangePriority(id domain.GoodsId, projectId domain.ProjectId) ([]*Good, error)
}

type goodServiceImpl struct {
	goodStorage GoodStorage
}

func NewGoodService(goodStorage GoodStorage) GoodService {
	return &goodServiceImpl{
		goodStorage: goodStorage,
	}
}

func (gs *goodServiceImpl) Create(projectId domain.ProjectId, name string) (*Good, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodServiceImpl) Update(id domain.GoodsId, projectId domain.ProjectId, name, description string) (*Good, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodServiceImpl) Delete(id domain.GoodsId, projectId domain.ProjectId) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodServiceImpl) List(limit, offset int) ([]*Good, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodServiceImpl) ChangePriority(id domain.GoodsId, projectId domain.ProjectId) ([]*Good, error) {
	//TODO implement me
	panic("implement me")
}
