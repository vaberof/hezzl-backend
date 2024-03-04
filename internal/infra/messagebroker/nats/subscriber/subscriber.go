package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vaberof/hezzl-backend/internal/infra/storage/clickhouse/chgoodlog"
)

const goodLogsSubject = "good.logs"

const defaultBatchSize = 10

type Subscriber interface {
	SubscribeOnGoodLogsSubject(goodLogStorage GoodLogStorage)
}

type subscriberImpl struct {
	natsConn *nats.Conn
}

func New(config *Config) (Subscriber, error) {
	nc, err := nats.Connect(fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return nil, err
	}
	return &subscriberImpl{natsConn: nc}, nil
}

func (s *subscriberImpl) SubscribeOnGoodLogsSubject(goodLogStorage GoodLogStorage) {
	goodLogs := make([]*GoodLog, 0, defaultBatchSize)

	s.natsConn.Subscribe(goodLogsSubject, func(msg *nats.Msg) {
		var goodLog GoodLog

		err := json.Unmarshal(msg.Data, &goodLog)
		if err != nil {
			return
		}

		goodLogs = append(goodLogs, &goodLog)

		if len(goodLogs) >= defaultBatchSize {
			err = goodLogStorage.Insert(buildCHGoodLogs(goodLogs))
			if err != nil {
				return
			}
			goodLogs = make([]*GoodLog, 0, defaultBatchSize)
		}
	})
}

func buildCHGoodLogs(goodLogs []*GoodLog) []*chgoodlog.GoodLog {
	chGoodLogs := make([]*chgoodlog.GoodLog, len(goodLogs))
	for i := range goodLogs {
		chGoodLogs[i] = buildCHGoodLog(goodLogs[i])
	}
	return chGoodLogs
}

func buildCHGoodLog(goodLog *GoodLog) *chgoodlog.GoodLog {
	return &chgoodlog.GoodLog{
		Id:          goodLog.Id,
		ProjectId:   goodLog.ProjectId,
		Name:        goodLog.Name,
		Description: goodLog.Description,
		Priority:    goodLog.Priority,
		Removed:     goodLog.Removed,
		EventTime:   goodLog.EventTime,
	}
}
