package domain

import (
	"context"
	"github.com/WLM1ke/gomoex"
)

// GroupUSD группа таблицы с курсом доллара.
const GroupUSD = "usd"

// USD - таблицы с курсом доллара.
type USD struct {
	ID

	iss *gomoex.ISSClient

	Rows []gomoex.Candle
}

func (U USD) Update(ctx context.Context) []Event {
	return nil
}

// UpdateUSD - правило обновления курса доллара.
//
// Если обновилась таблица с торговыми датами, то нужно обновить курс.
type UpdateUSD struct {
	iss *gomoex.ISSClient
}

// NewUpdateUSD - создает новое правило обновления доллара.
func NewUpdateUSD(iss *gomoex.ISSClient) *UpdateUSD {
	return &UpdateUSD{iss: iss}
}

// Activate активирует правило.
func (u *UpdateUSD) Activate(ctx context.Context, in <-chan Event, out chan<- Event) {
	for {
		select {
		case event := <-in:
			if event.Group() != GroupTradingDates || event.Name() != GroupTradingDates {
				continue
			}

			_, ok := event.(RowsReplaced)
			if !ok {
				continue
			}

			out <- UpdateRequired{&USD{ID: NewID(GroupUSD, GroupUSD), iss: u.iss}}
		case <-ctx.Done():
			return
		}
	}
}
