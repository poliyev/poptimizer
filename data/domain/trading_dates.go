package domain

import (
	"context"
	"time"

	"github.com/WLM1ke/gomoex"
	"go.uber.org/zap"
)

// GroupTradingDates - группа таблицы с торговыми данными.
const GroupTradingDates = "trading_dates"

// Информация о торгах публикуется на MOEX ISS в 0:45 по московскому времени на следующий день.
const (
	issTZ     = "Europe/Moscow"
	issHour   = 0
	issMinute = 45
)

func prepareZone(tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		zap.L().Panic("Не удалось загрузить часовой пояс", zap.Error(err))
	}

	return loc
}

func nextISSDailyUpdate(now time.Time, tz *time.Location) time.Time {
	now = now.In(tz)
	end := time.Date(now.Year(), now.Month(), now.Day(), issHour, issMinute, 0, 0, tz)

	if end.Before(now) {
		end = end.AddDate(0, 0, 1)
	}

	return end
}

// CheckTradingDay формирует команды о необходимости проверки окончания торгового дня.
//
// Требуется при запуске приложения и ежедневно после публикации данных на MOEX ISS.
type CheckTradingDay struct {
	ticker <-chan time.Time
	stopFn func()
	tz     *time.Location
}

// NewCheckTradingDay создает источник сообщений о начале торгового дня.
//
// Так как компьютер может заснуть, что вызывает расхождение между монотонным и фактическим временем,
// то проверку публикации данных лучше проводить на регулярной основе, а не привязать к конкретному времени.
func NewCheckTradingDay() *CheckTradingDay {
	ticker := time.NewTicker(time.Hour)

	return &CheckTradingDay{
		ticker: ticker.C,
		stopFn: ticker.Stop,
		tz:     prepareZone(issTZ),
	}
}

// StartProduceCommands записывает команду о необходимости проверки обновления таблицы с торговыми датами на регулярной основе.
func (d *CheckTradingDay) StartProduceCommands(ctx context.Context, output chan<- Command) {
	cmd := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}
	output <- &cmd

	now := time.Now()
	nextDataUpdate := nextISSDailyUpdate(now, d.tz)

	for {
		select {
		case now = <-d.ticker:
			if now.After(nextDataUpdate) {
				output <- &cmd

				nextDataUpdate = nextISSDailyUpdate(now, d.tz)
			}
		case <-ctx.Done():
			d.stopFn()

			return
		}
	}
}

// TradingDates - таблица с диапазоном торговых дат.
type TradingDates struct {
	TableID

	iss *gomoex.ISSClient

	Rows []gomoex.Date
}

// HandleCommand - полностью переписывает таблицу, если появились новые данные о торговых датах.
func (t *TradingDates) HandleCommand(ctx context.Context, _ Command) []Event {
	newRows, err := t.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)

	switch {
	case err != nil:
		zap.L().Panic("Не удалось получить данные ISS", zap.Error(err))
	case len(newRows) != 1:
		zap.L().Panic("Ошибка валидации данных ISS", zap.Error(err))
	case t.Rows == nil, !newRows[0].Till.Equal(t.Rows[0].Till):
		return []Event{RowsReplaced{t.TableID, newRows}}
	}

	return nil
}
