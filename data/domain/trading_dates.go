package domain

import (
	"context"
	"time"

	"github.com/WLM1ke/gomoex"
	"go.uber.org/zap"
)

// GroupTradingDates - группа таблицы с торговыми данными.
const GroupTradingDates = "trading_dates"

func prepareZone(zone string) *time.Location {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		zap.L().Panic("Не удалось загрузить часовой пояс Москвы", zap.Error(err))
	}

	return loc
}

// Информация о торгах публикуется на MOEX ISS в 0:45 по московскому времени на следующий день.
var issZone = prepareZone("Europe/Moscow")

const (
	issHour   = 0
	issMinute = 45
)

func nextISSDailyUpdate(now time.Time) time.Time {
	now = now.In(issZone)
	end := time.Date(now.Year(), now.Month(), now.Day(), issHour, issMinute, 0, 0, issZone)

	if end.Before(now) {
		end = end.AddDate(0, 0, 1)
	}

	return end
}

// CheckTradingDay формирует команды о необходимости проверки окончания торгового дня.
//
// Требуется при запуске приложения и ежедневно после публикации данных на MOEX ISS.
type CheckTradingDay struct {
	timer <-chan time.Time
}

// NewCheckTradingDay создает источник сообщений о начале торгового дня.
//
// Так как компьютер может заснуть, что вызывает расхождение между монотонным и фактическим временем,
// то проверку публикации данных лучше проводить на регулярной основе, а не привязать к конкретному времени.
func NewCheckTradingDay() *CheckTradingDay {
	return &CheckTradingDay{time.Tick(time.Hour)}
}

// StartProduceCommands записывает команду о необходимости проверки обновления таблицы с торговыми датами на регулярной основе.
func (d CheckTradingDay) StartProduceCommands(ctx context.Context, output chan<- Command) {
	timer := d.timer

	cmd := UpdateTable{TableID{GroupTradingDates, GroupTradingDates}}
	output <- &cmd

	now := time.Now()
	nextDataUpdate := nextISSDailyUpdate(now)

	for {
		select {
		case now = <-timer:
			if now.After(nextDataUpdate) {
				output <- &cmd
				nextDataUpdate = nextISSDailyUpdate(now)
			}
		case <-ctx.Done():
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
