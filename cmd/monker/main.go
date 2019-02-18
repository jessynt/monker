package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"monker/pkg/config"
	"monker/pkg/consumer"
	nsqLogger "monker/pkg/nsq-logger"
	"monker/pkg/producer"
	"monker/pkg/storage"
	"monker/pkg/workers"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/youzan/go-nsq"
)

var (
	Version   = "unknown"
	BuildDate = "unknown"
)

type startFunc func(chan error)

func main() {
	app := &cli.App{
		Name:    "Monker",
		Version: fmt.Sprintf("%s+%s", Version, BuildDate),
		Action:  run,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config,c",
				Usage:  "Source of a configuration file",
				Value:  "config.yaml",
				EnvVar: "MONKER_CONFIG_PATH",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Panic(err)
	}
}

func run(ctx *cli.Context) error {
	log.WithField("version", Version).Info("Monker starting...")

	globalConfig, err := config.Load(ctx.String("config"))
	if err != nil {
		return err
	}
	logLevel, err := log.ParseLevel(globalConfig.LogLevel)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	consumerConfig := nsq.NewConfig()
	consumerConfig.EnableOrdered = true
	nsqConsumer, err := nsq.NewConsumer(globalConfig.Nsq.Topic, globalConfig.Nsq.Channel, consumerConfig)
	if err != nil {
		return err
	}

	nsqConsumer.SetLogger(nsqLogger.NewNSQLogrusLoggerAtLevel(logLevel))

	kafkaWrite := kafka.NewWriter(kafka.WriterConfig{
		Brokers: globalConfig.Kafka.Brokers,
		Topic:   globalConfig.Kafka.Topic,
	})

	kafkaProducer, err := producer.NewKafkaProducer(kafkaWrite)
	if err != nil {
		return err
	}

	storageURL, err := url.Parse(globalConfig.StorageDSN)
	if err != nil {
		return err
	}

	s, err := storage.NewStorage(storageURL)
	if err != nil {
		return err
	}

	monkerWorker, err := workers.NewMonkerWorker(globalConfig.Worker, kafkaProducer, s)
	if err != nil {
		return err
	}

	c, err := consumer.NewConsumer(globalConfig.Nsq, nsqConsumer, monkerWorker.MessageHandler)
	if err != nil {
		return err
	}

	starts := []startFunc{
		func(errs chan error) {
			if err := c.Start(); err != nil {
				errs <- err
			}
		},
		func(errs chan error) {
			errs <- monkerWorker.Start()
		},
		func(errs chan error) {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGABRT, os.Kill)

			s := <-c
			log.WithField("signal", s.String()).Info("caught signal")
			errs <- nil
		},
	}

	errs := make(chan error, len(starts))
	for _, c := range starts {
		go c(errs)
	}

	select {
	case err := <-errs:
		if err := c.Stop(); err != nil {
			log.WithError(err).Error("failed to stop consumer")
		}
		if err := monkerWorker.Stop(); err != nil {
			log.WithError(err).Error("failed to stop monker worker")
		}
		return err
	}
}
