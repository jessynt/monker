package consumer

import (
	"monker/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/youzan/go-nsq"
)

type Consumer struct {
	config      config.NsqConfig
	nsqConsumer *nsq.Consumer
}

func NewConsumer(c config.NsqConfig, nsqConsumer *nsq.Consumer, messageHandler MessageHandler) (*Consumer, error) {
	nsqConsumer.AddConcurrentHandlers(&ExchangeHandler{messageHandler: messageHandler}, c.Concurrency)
	return &Consumer{
		config:      c,
		nsqConsumer: nsqConsumer,
	}, nil
}

func (c *Consumer) Start() error {
	if err := c.nsqConsumer.ConnectToNSQLookupds(c.config.Lookupds); err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Stop() error {
	log.Info("consumer stopping...")
	c.nsqConsumer.Stop()
	<-c.nsqConsumer.StopChan // block until stop completed
	log.Info("consume stopped")
	return nil
}
