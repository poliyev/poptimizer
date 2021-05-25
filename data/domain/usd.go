package domain

import (
	"context"
	"fmt"
	"github.com/WLM1ke/gomoex"
	"time"
)

// GroupUSD группа таблицы с курсом доллара.
const GroupUSD = "usd"

// CandlesGateway получает информацию о свечках.
type CandlesGateway interface {
	TradingDatesGateway
	MarketCandles(ctx context.Context, engine, market, ticker, start, last string, interval int) ([]gomoex.Candle, error)
}

func dateToISSString(date time.Time) string {
	return date.Format("2006-01-02")
}

// USD - таблицы с курсом доллара.
type USD struct {
	ID

	iss   CandlesGateway
	dates *TradingDates

	Rows []gomoex.Candle
}

func (u USD) Update(ctx context.Context) []Event {
	startDate := ""
	if u.Rows != nil {
		startDate = dateToISSString(u.Rows[len(u.Rows)-1].Begin)
	}

	lastDate := dateToISSString(u.dates.lastDate())

	newData, err := u.iss.MarketCandles(ctx, "currency", "selt", "USD000UTSTOM", startDate, lastDate, gomoex.IntervalDay)

	switch {
	case err != nil:
		return []Event{UpdateError{u.ID, err}}
	case u.Rows == nil:
		return []Event{RowsAppended{u.ID, newData}}
	case u.Rows[len(u.Rows)-1] != newData[0]:
		err = fmt.Errorf("не совпадают данные %+v != %+v", u.Rows[len(u.Rows)-1], newData[0])

		return []Event{UpdateError{u.ID, err}}
	default:
		return []Event{RowsAppended{u.ID, newData[1:]}}
	}
}

// UpdateUSD - правило обновления курса доллара.
//
// Если обновилась таблица с торговыми датами, то нужно обновить курс.
type UpdateUSD struct {
	iss CandlesGateway
}

// NewUpdateUSD - создает новое правило обновления доллара.
func NewUpdateUSD(iss CandlesGateway) *UpdateUSD {
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

			out <- makeUSDEvent(u.iss)
		case <-ctx.Done():
			return
		}
	}
}

func makeUSDEvent(iss CandlesGateway) UpdateRequired {
	dates := TradingDates{ID: ID{GroupTradingDates, GroupTradingDates}, iss: iss}

	return UpdateRequired{
		[]Table{
			&USD{ID: NewID(GroupUSD, GroupUSD), iss: iss, dates: &dates},
			&dates,
		},
	}
}
