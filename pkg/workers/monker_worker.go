package workers

import (
	"context"
	"sync"
	"time"

	"monker/pkg/config"
	"monker/pkg/producer"
	"monker/pkg/storage"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	errPutToStorage = errors.New("failed to put message to storage")
)

type MonkerWorker struct {
	sync.Mutex
	config            config.WorkerConfig
	producer          producer.Producer
	storage           storage.Storage
	cache             []*producer.Message
	lastFlushAt       time.Time
	readStorageTicker *time.Ticker
	stopChan          chan chan struct{}
}

func NewMonkerWorker(
	workerConfig config.WorkerConfig,
	p producer.Producer,
	s storage.Storage,
) (*MonkerWorker, error) {
	return &MonkerWorker{
		config:   workerConfig,
		producer: p,
		storage:  s,
		stopChan: make(chan chan struct{}),
	}, nil
}

func (w *MonkerWorker) Execute(ctx context.Context) {
	w.Lock()
	defer w.Unlock()

	if len(w.cache) >= w.config.CacheSize || time.Now().Sub(w.lastFlushAt) >= w.config.CacheFlushTimeout {
		log.WithFields(log.Fields{"len": len(w.cache), "last_flush_at": w.lastFlushAt}).
			Debug("Flushing worker cache to Kafka")

		if len(w.cache) > 0 {
			messages := make([]*producer.Message, len(w.cache))
			copy(messages, w.cache)
			w.cache = w.cache[:0] // will keep allocated memory

			go w.publishMessages(ctx, messages)
		}
		w.lastFlushAt = time.Now()
	}
}

func (w *MonkerWorker) Start() error {
	w.readStorageTicker = time.NewTicker(w.config.StorageReadTimeout)
	for {
		select {
		case stopFinished := <-w.stopChan:
			stopFinished <- struct{}{}
			return nil
		case <-w.readStorageTicker.C:
			log.Debug("Monker worker Tick")
			w.populateCacheFromStorage()
		default:
			w.Execute(context.Background())
		}

		log.WithField("timeout", w.config.CycleTimeout.String()).
			Debug("Monker worker is going to sleep for a while")
		time.Sleep(w.config.CycleTimeout)
	}
}

func (w *MonkerWorker) publishMessages(ctx context.Context, messages []*producer.Message) {
	for _, message := range messages {
		err := w.producer.Publish(ctx, *message)
		if err != nil {
			log.WithError(err).WithField("msg", message.String()).
				Warning("Failed to publish messages to Kafka, moving to storage")

			if err = w.storeMessage(message); err != nil {
				if err == errPutToStorage {
					w.cacheMessage(message)
				} else {
					log.WithError(err).WithField("msg", message.String()).Error("Unhandled storage error")
				}
			}
		}
	}
}

func (w *MonkerWorker) cacheMessage(message *producer.Message) error {
	w.Lock()
	defer w.Unlock()

	w.cache = append(w.cache, message)
	return nil
}

func (w *MonkerWorker) MessageHandler(body []byte) error {
	return w.cacheMessage(producer.NewMessage(body))
}

func (w *MonkerWorker) storeMessage(message *producer.Message) error {
	err := w.storage.Put(message.Body)
	if err != nil {
		return errPutToStorage
	}
	return nil
}

func (w *MonkerWorker) Stop() error {
	log.Info("monker worker stopping")

	stopDone := make(chan struct{}, 1)
	w.stopChan <- stopDone
	w.readStorageTicker.Stop()
	<-stopDone

	return w.storage.Close()
}

func (w *MonkerWorker) populateCacheFromStorage() {
	log.Debug("Populating cache from storage")
	for {
		message, err := w.storage.Get()
		if err != nil {
			if err == storage.ErrStorageIsEmpty {
				break
			}
			log.WithError(err).Error("Failed to read message from persistent storage")
			continue
		}
		w.cacheMessage(producer.NewMessage(message))
	}
}
