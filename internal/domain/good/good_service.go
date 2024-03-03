package good

import (
	"encoding/json"
	"errors"
	"github.com/vaberof/hezzl-backend/internal/infra/storage"
	"github.com/vaberof/hezzl-backend/pkg/domain"
	"strconv"
	"time"
)

const (
	goodKey     = "good_"
	goodListKey = "good_list_"
	limitKey    = "limit_"
	offsetKey   = "offset_"
)

const (
	goodListCacheExpireTime = 1 * time.Minute
)

var (
	ErrGoodNotFound = errors.New("good not found")
)

type GoodService interface {
	Create(projectId domain.ProjectId, name domain.GoodName) (*Good, error)
	Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*Good, error)
	Delete(id domain.GoodId, projectId domain.ProjectId) error
	List(limit, offset int) ([]*Good, error)
	ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority int) ([]*Good, error)
}

type goodServiceImpl struct {
	goodStorage     GoodStorage
	inMemoryStorage InMemoryStorage
}

func NewGoodService(goodStorage GoodStorage, inMemoryStorage InMemoryStorage) GoodService {
	return &goodServiceImpl{
		goodStorage:     goodStorage,
		inMemoryStorage: inMemoryStorage,
	}
}

func (g *goodServiceImpl) Create(projectId domain.ProjectId, name domain.GoodName) (*Good, error) {
	// log to clickhouse

	domainGood, err := g.goodStorage.Create(projectId, name)
	if err != nil {
		return nil, err
	}

	return domainGood, nil
}

func (g *goodServiceImpl) Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description *domain.GoodDescription) (*Good, error) {
	// log to clickhouse

	exists, err := g.goodStorage.IsExists(id, projectId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrGoodNotFound
	}

	domainGood, err := g.goodStorage.Update(id, projectId, name, description)
	if err != nil {
		return nil, err
	}

	goodCacheKey := g.getGoodCacheKey(id, projectId)

	err = g.inMemoryStorage.Delete(goodCacheKey)
	if err != nil {
		if !errors.Is(err, storage.ErrRedisKeyNotFound) {
			return nil, err
		}
	}

	return domainGood, nil
}

func (g *goodServiceImpl) Delete(id domain.GoodId, projectId domain.ProjectId) error {
	// log to clickhouse

	err := g.goodStorage.Delete(id, projectId)
	if err != nil {
		return err
	}

	goodCacheKey := g.getGoodCacheKey(id, projectId)

	err = g.inMemoryStorage.Delete(goodCacheKey)
	if err != nil {
		if !errors.Is(err, storage.ErrRedisKeyNotFound) {
			return err
		}
	}

	return nil
}

func (g *goodServiceImpl) List(limit, offset int) ([]*Good, error) {
	// log to clickhouse

	goodListCacheKey := g.getGoodListCacheKey(limit, offset)

	cachedDomainGoods, err := g.getCachedGoods(goodListCacheKey)
	if err == nil {
		return cachedDomainGoods, nil
	}
	if err != nil {
		if !errors.Is(err, storage.ErrRedisKeyNotFound) {
			return nil, err
		}
	}

	domainGoods, err := g.goodStorage.List(limit, offset)
	if err != nil {
		return nil, err
	}

	domainGoodsBytes, err := json.Marshal(&domainGoods)
	if err != nil {
		return nil, err
	}

	err = g.inMemoryStorage.Set(goodListCacheKey, string(domainGoodsBytes), goodListCacheExpireTime)
	if err != nil {
		return nil, err
	}

	return domainGoods, nil
}

func (g *goodServiceImpl) ChangePriority(id domain.GoodId, projectId domain.ProjectId, newPriority int) ([]*Good, error) {
	// log to clickhouse

	exists, err := g.goodStorage.IsExists(id, projectId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrGoodNotFound
	}

	domainGoods, err := g.goodStorage.ChangePriority(id, projectId, newPriority)
	if err != nil {
		return nil, err
	}

	goodCacheKeys := g.getGoodCacheKeys(domainGoods)

	err = g.inMemoryStorage.Delete(goodCacheKeys...)
	if !errors.Is(err, storage.ErrRedisKeyNotFound) {
		return nil, err
	}

	return domainGoods, nil
}

func (g *goodServiceImpl) getCachedGoods(key string) ([]*Good, error) {
	cachedGoodsStr, err := g.inMemoryStorage.Get(key)
	if err != nil {
		return nil, err
	}

	var domainGoods []*Good

	err = json.Unmarshal([]byte(cachedGoodsStr), &domainGoods)
	if err != nil {
		return nil, err
	}

	return domainGoods, nil
}

func (g *goodServiceImpl) getGoodCacheKeys(domainGoods []*Good) []string {
	goodCacheKeys := make([]string, len(domainGoods))
	for i := range domainGoods {
		goodCacheKeys[i] = g.getGoodCacheKey(domainGoods[i].Id, domainGoods[i].ProjectId)
	}
	return goodCacheKeys
}

func (g *goodServiceImpl) getGoodCacheKey(id domain.GoodId, projectId domain.ProjectId) string {
	idStr := strconv.Itoa(int(id))
	projectIdStr := strconv.Itoa(int(projectId))
	goodCacheKey := goodKey + idStr + "_" + projectIdStr
	return goodCacheKey
}

func (g *goodServiceImpl) getGoodListCacheKey(limit, offset int) string {
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	goodListCacheKey := goodListKey + limitKey + limitStr + "_" + offsetKey + offsetStr
	return goodListCacheKey
}
