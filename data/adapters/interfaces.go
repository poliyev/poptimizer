package adapters

import (
	"context"
	"poptimizer/data/domain"
)

// TableRepo реализует загрузку и сохранение изменения таблиц.
type TableRepo interface {
	Load(ctx context.Context, id domain.TableID) (domain.Table, error)
	Save(ctx context.Context, event domain.Event) error
}

// JSONViewer обеспечивает просмотр данных таблиц в формате ExtendedJSON.
type JSONViewer interface {
	ViewJSON(ctx context.Context, id domain.TableID) ([]byte, error)
}
