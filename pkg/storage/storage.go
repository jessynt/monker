package storage

import (
	"net/url"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUnknownStorage = errors.New("unknown storage type")
	ErrStorageIsEmpty = errors.New("No more messages in storage")
)

type Storage interface {
	Put(data []byte) error
	Get() ([]byte, error)
	Close() error
}

func NewStorage(dsn *url.URL) (Storage, error) {
	log.WithField("dsn", dsn.String()).Debug("Trying to instantiate new storage instance")

	log.WithField("type", dsn.Scheme).Info("Looking for storage")
	switch dsn.Scheme {
	case "inmem":
		return NewInmemStorage(), nil
	}
	return nil, ErrUnknownStorage
}
