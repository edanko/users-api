package commands

import (
	"context"
)

type eventBus interface {
	Publish(ctx context.Context, event any) error
}
