package workers

import (
	"context"
)

type Worker interface {
	Execute(ctx context.Context)
	Start() error
	Stop() error
}
