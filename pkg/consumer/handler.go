package consumer

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/youzan/go-nsq"
)

type MessageHandler func(body []byte) error

type ExchangeHandler struct {
	messageHandler MessageHandler
}

func (h *ExchangeHandler) HandleMessage(message *nsq.Message) error {
	logrus.WithField("message", string(message.Body)).Debug("consuming message...")
	err := h.messageHandler(message.Body)
	if err != nil {
		message.RequeueWithoutBackoff(time.Duration(5) * time.Second)
		return err
	}
	message.Finish()
	return nil
}
