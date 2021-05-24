package app

import (
	"context"
	"go.uber.org/zap"
	"poptimizer/data/adapters"
	"poptimizer/data/domain"
)

// ErrorsHandler - правило обработки ошибок.
type ErrorsHandler struct {
}

// Activate - активирует правило.
//
// Реагирует паникой на событие ошибок и не использует исходящий канал.
func (e ErrorsHandler) Activate(ctx context.Context, in <-chan domain.Event, _ chan<- domain.Event) {
	for {
		select {
		case event := <-in:
			errorEvent, ok := event.(domain.UpdateError)
			if !ok {
				continue
			}
			zap.L().Panic("Error", adapters.EventField(errorEvent), zap.Error(errorEvent.Error))
		case <-ctx.Done():
			return
		}
	}
}
