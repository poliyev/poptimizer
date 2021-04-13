package cfg

import "context"

type Module interface {
	Name() string
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
