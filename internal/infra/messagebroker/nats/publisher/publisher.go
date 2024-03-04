package publisher

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

const goodLogsSubject = "good.logs"

type Publisher interface {
	PublishGoodLog(id, projectId int64, name, description string, priority int, removed bool, eventTime time.Time) error
}

type publisherImpl struct {
	natsConn *nats.Conn
}

func New(config *Config) (Publisher, error) {
	nc, err := nats.Connect(fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return nil, err
	}
	return &publisherImpl{natsConn: nc}, nil
}

func (p *publisherImpl) PublishGoodLog(id, projectId int64, name, description string, priority int, removed bool, eventTime time.Time) error {
	data, err := json.Marshal(&GoodLog{
		Id:          id,
		ProjectId:   projectId,
		Name:        name,
		Description: description,
		Priority:    priority,
		Removed:     removed,
		EventTime:   eventTime,
	})
	if err != nil {
		return err
	}
	err = p.natsConn.Publish(goodLogsSubject, data)
	if err != nil {
		return err
	}
	return nil
}

func (p *publisherImpl) Publish(subject string, data []byte) error {
	err := p.natsConn.Publish(subject, data)
	if err != nil {
		return err
	}
	return nil
}
