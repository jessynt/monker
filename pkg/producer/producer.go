package producer

import (
	"context"
)

type Producer interface {
	Publish(ctx context.Context, message Message) error
	Close() error
}
