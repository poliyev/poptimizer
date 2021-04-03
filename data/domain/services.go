package domain

import (
	"context"
	"time"
)

func prepareLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("не удалось загрузить часовой пояс Москвы")
	}

	return loc
}

var zoneMoscow = prepareLocation()

//Информация о торгах публикуется на MOEX ISS в 0:45 на следующий день.
func nextISSDailyUpdate() time.Time {
	now := time.Now().In(zoneMoscow)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, zoneMoscow)

	if end.Before(now) {
		end = end.AddDate(0, 0, 1)
	}

	return end
}

// CheckTradingDay формирует команды о необходимости проверки окончания торгового дня.
//
// Требуется при запуске приложения и ежедневно после публикации данных на MOEX ISS.
type CheckTradingDay struct {
}

func (d CheckTradingDay) StartProduceCommands(ctx context.Context, output chan<- Command) {
	cmd := UpdateTable{ID{GroupTradingDates, GroupTradingDates}}

	output <- &cmd
	nextDataUpdate := nextISSDailyUpdate()

	// Так как компьютер может заснуть, что вызывает расхождение между монотонным фактическим временем,
	// то проверку публикации данных лучше проводить на регулярной основе, а не привязать к конкретному времени.
	timer := time.Tick(time.Hour)

LOOP:
	for {
		select {
		case <-timer:
			if time.Now().After(nextDataUpdate) {
				output <- &cmd
				nextDataUpdate = nextISSDailyUpdate()
			}
		case <-ctx.Done():
			break LOOP
		}
	}
}
