package pggood

import (
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/pkg/domain"
)

func toDomainGoods(postgresGoods []*Good) []*good.Good {
	domainGoods := make([]*good.Good, len(postgresGoods))
	for i := range postgresGoods {
		domainGoods[i] = toDomainGood(postgresGoods[i])
	}
	return domainGoods
}

func toDomainGood(postgresGood *Good) *good.Good {
	return &good.Good{
		Id:          domain.GoodId(postgresGood.Id),
		ProjectId:   domain.ProjectId(postgresGood.ProjectId),
		Name:        domain.GoodName(postgresGood.Name),
		Description: domain.GoodDescription(postgresGood.Description),
		Priority:    domain.GoodPriority(postgresGood.Priority),
		Removed:     domain.GoodRemoved(postgresGood.Removed),
		CreatedAt:   domain.GoodCreatedAt(postgresGood.CreatedAt),
	}
}
