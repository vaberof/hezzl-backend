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
	Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description domain.GoodDescription) (*Good, error)
	Delete(id domain.GoodId, projectId domain.ProjectId) error
	List(limit, offset int) ([]*Good, error)
	ChangePriority(id domain.GoodId, projectId domain.ProjectId) ([]*Good, error)
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

func (gs *goodServiceImpl) Create(projectId domain.ProjectId, name domain.GoodName) (*Good, error) {
	// log to clickhouse

	domainGood, err := gs.goodStorage.Create(projectId, name)
	if err != nil {
		return nil, err
	}

	return domainGood, nil
}

func (gs *goodServiceImpl) Update(id domain.GoodId, projectId domain.ProjectId, name domain.GoodName, description domain.GoodDescription) (*Good, error) {
	// log to clickhouse

	exists, err := gs.goodStorage.IsExists(id, projectId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrGoodNotFound
	}

	domainGood, err := gs.Update(id, projectId, name, description)
	if err != nil {
		return nil, err
	}

	goodCacheKey := gs.getGoodCacheKey(id, projectId)

	err = gs.inMemoryStorage.Delete(goodCacheKey)
	if err != nil {
		if !errors.Is(err, storage.ErrRedisKeyNotFound) {
			return nil, err
		}
	}

	return domainGood, nil
}

func (gs *goodServiceImpl) Delete(id domain.GoodId, projectId domain.ProjectId) error {
	// log to clickhouse

	err := gs.goodStorage.Delete(id, projectId)
	if err != nil {
		return err
	}

	goodCacheKey := gs.getGoodCacheKey(id, projectId)

	err = gs.inMemoryStorage.Delete(goodCacheKey)
	if err != nil {
		if !errors.Is(err, storage.ErrRedisKeyNotFound) {
			return err
		}
	}

	return nil
}

func (gs *goodServiceImpl) List(limit, offset int) ([]*Good, error) {
	// log to clickhouse

	goodListCacheKey := gs.getGoodListCacheKey(limit, offset)

	cachedDomainGoods, err := gs.getCachedGoods(goodListCacheKey)
	if err == nil {
		return cachedDomainGoods, nil
	}
	if err != nil {
		if !errors.Is(err, storage.ErrRedisKeyNotFound) {
			return nil, err
		}
	}

	domainGoods, err := gs.goodStorage.List(limit, offset)
	if err != nil {
		return nil, err
	}

	domainGoodsBytes, err := json.Marshal(&domainGoods)
	if err != nil {
		return nil, err
	}

	err = gs.inMemoryStorage.Set(goodListCacheKey, string(domainGoodsBytes), goodListCacheExpireTime)
	if err != nil {
		return nil, err
	}

	return domainGoods, nil
}

func (gs *goodServiceImpl) ChangePriority(id domain.GoodId, projectId domain.ProjectId) ([]*Good, error) {
	// log to clickhouse

	exists, err := gs.goodStorage.IsExists(id, projectId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrGoodNotFound
	}

	domainGoods, err := gs.goodStorage.ChangePriority(id, projectId)
	if err != nil {
		return nil, err
	}

	goodCacheKeys := gs.getGoodCacheKeys(domainGoods)

	err = gs.inMemoryStorage.Delete(goodCacheKeys...)
	if !errors.Is(err, storage.ErrRedisKeyNotFound) {
		return nil, err
	}

	return domainGoods, nil
}

func (gs *goodServiceImpl) getCachedGoods(key string) ([]*Good, error) {
	cachedGoodsStr, err := gs.inMemoryStorage.Get(key)
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

func (gs *goodServiceImpl) getGoodCacheKeys(domainGoods []*Good) []string {
	goodCacheKeys := make([]string, len(domainGoods))
	for i := range domainGoods {
		goodCacheKeys[i] = gs.getGoodCacheKey(domainGoods[i].Id, domainGoods[i].ProjectId)
	}
	return goodCacheKeys
}

func (gs *goodServiceImpl) getGoodCacheKey(id domain.GoodId, projectId domain.ProjectId) string {
	idStr := strconv.Itoa(int(id))
	projectIdStr := strconv.Itoa(int(projectId))
	goodCacheKey := goodKey + idStr + "_" + projectIdStr
	return goodCacheKey
}

func (gs *goodServiceImpl) getGoodListCacheKey(limit, offset int) string {
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	goodListCacheKey := goodListKey + limitKey + limitStr + "_" + offsetKey + offsetStr
	return goodListCacheKey
}
