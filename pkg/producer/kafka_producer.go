package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(write *kafka.Writer) (Producer, error) {
	return &KafkaProducer{
		writer: write,
	}, nil
}

func (p *KafkaProducer) Publish(ctx context.Context, message Message) (err error) {
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   message.ID.Bytes(),
		Value: message.Body,
	})

	if err == nil {
		log.WithField("message", message.String()).Debug("Successfully sent message to kafka")
	} else {
		log.WithError(err).WithField("message", message.String()).Error("Failed to publish message to kafka")
	}

	return err
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
