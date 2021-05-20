package app

import (
	"context"
	"poptimizer/data/domain"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Repo осуществляет восстановление талицы и сохранение новых строк.
type Repo interface {
	Unmarshal(ctx context.Context, event domain.UpdateRequired) (domain.Table, error)
	Replace(ctx context.Context, event domain.RowsReplaced) error
	Append(ctx context.Context, event domain.RowsAppended) error
}

// UoW контролирует процесс загрузки существующих данных в таблицу, ее обновления и сохранения изменений.
//
// Один цикл работы с таблицей должен укладываться в timeout и реализует интерфейс бизнес-правила.
type UoW struct {
	repo    Repo
	timeout time.Duration
}

// NewUoW создает Unit of Work.
func NewUoW(repo Repo, timeout time.Duration) *UoW {
	return &UoW{repo, timeout}
}

// Activate запускает работу бизнес правила.
//
// Выбирает события обновления таблиц, сохраняет события с изменениями в репозиторий и посылает их в исходящий канал.
// Таким образом, события связанные с обновлением таблиц гарантировано поступают на обработку другими правилами после сохранения результатов
// в репозитории.
func (u UoW) Activate(ctx context.Context, in <-chan domain.Event, out chan<- domain.Event) {
	wg := sync.WaitGroup{}

	for {
		select {
		case event := <-in:
			update, ok := event.(domain.UpdateRequired)
			if !ok {
				continue
			}

			wg.Add(1)

			go func() {
				defer wg.Done()

				eventCtx, cancel := context.WithTimeout(ctx, u.timeout)
				defer cancel()

				table := u.unmarshalTable(eventCtx, update)

				for _, event := range table.Update(eventCtx) {
					u.saveChanges(eventCtx, event)
					out <- event
				}
			}()

		case <-ctx.Done():
			wg.Wait()

			return
		}
	}
}

func (u UoW) unmarshalTable(ctx context.Context, event domain.UpdateRequired) domain.Table {
	table, err := u.repo.Unmarshal(ctx, event)
	if err != nil {
		zap.L().Panic("Unmarshal table", zap.Stringer("id", event), zap.Error(err))
	}

	return table
}

func (u UoW) saveChanges(ctx context.Context, event domain.Event) {
	switch event := event.(type) {
	case domain.RowsReplaced:
		if err := u.repo.Replace(ctx, event); err != nil {
			zap.L().Panic("Replace rows", zap.Stringer("id", event), zap.Error(err))
		}
	case domain.RowsAppended:
		if err := u.repo.Append(ctx, event); err != nil {
			zap.L().Panic("Append rows", zap.Stringer("id", event), zap.Error(err))
		}
	}
}
