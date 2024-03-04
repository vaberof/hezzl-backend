package subscriber

import "github.com/vaberof/hezzl-backend/internal/infra/storage/clickhouse/chgoodlog"

type GoodLogStorage interface {
	Insert(goodLogs []*chgoodlog.GoodLog) error
}
