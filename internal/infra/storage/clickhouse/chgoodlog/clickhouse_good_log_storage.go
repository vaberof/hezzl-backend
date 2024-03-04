package chgoodlog

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"log"
)

type ClickHouseGoodLogStorage struct {
	chConn clickhouse.Conn
}

func NewCHGoodLogStorage(chConn clickhouse.Conn) *ClickHouseGoodLogStorage {
	return &ClickHouseGoodLogStorage{chConn: chConn}
}

func (ch *ClickHouseGoodLogStorage) Insert(goodLogs []*GoodLog) error {
	query := `
		INSERT INTO good_logs
	`

	batch, err := ch.chConn.PrepareBatch(context.Background(), query)
	if err != nil {
		return err
	}

	for i := range goodLogs {
		err = batch.Append(
			&goodLogs[i].Id,
			goodLogs[i].ProjectId,
			goodLogs[i].Name,
			goodLogs[i].Description,
			goodLogs[i].Priority,
			goodLogs[i].Removed,
			goodLogs[i].EventTime,
		)
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		log.Println("failed to send a batch to clickhouse", err)
	} else {
		log.Println("sent a batch to clickhouse successfully", err)
	}

	return err
}
