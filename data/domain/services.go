package domain

import (
	"context"
	"time"
)

func prepareZone() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("Не удалось загрузить часовой пояс Москвы")
	}

	return loc
}

//Информация о торгах публикуется на MOEX ISS в 0:45 по московскому времени на следующий день.
var zoneMoscow = prepareZone()

const (
	issHour   = 0
	issMinute = 45
)

func nextISSDailyUpdate(now time.Time) time.Time {
	now = now.In(zoneMoscow)
	end := time.Date(now.Year(), now.Month(), now.Day(), issHour, issMinute, 0, 0, zoneMoscow)

	if end.Before(now) {
		end = end.AddDate(0, 0, 1)
	}

	return end
}

// Так как компьютер может заснуть, что вызывает расхождение между монотонным и фактическим временем,
// то проверку публикации данных лучше проводить на регулярной основе, а не привязать к конкретному времени.
var defaultUpdateTimer = time.Tick(time.Hour)

// CheckTradingDay формирует команды о необходимости проверки окончания торгового дня.
//
// Требуется при запуске приложения и ежедневно после публикации данных на MOEX ISS.
type CheckTradingDay struct {
	timer <-chan time.Time
}

func (d CheckTradingDay) StartProduceCommands(ctx context.Context, output chan<- Command) {
	timer := d.timer
	if timer == nil {
		timer = defaultUpdateTimer
	}

	cmd := UpdateTable{ID{GroupTradingDates, GroupTradingDates}}

	now := time.Now()
	output <- &cmd
	nextDataUpdate := nextISSDailyUpdate(now)

LOOP:
	for {
		select {
		case now = <-timer:
			if now.After(nextDataUpdate) {
				output <- &cmd
				nextDataUpdate = nextISSDailyUpdate(now)
			}
		case <-ctx.Done():
			break LOOP
		}
	}
}
