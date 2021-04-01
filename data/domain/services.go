package domain

import (
	"context"
	"time"
)

type WorkStarted struct {
}

func (w WorkStarted) StartProduceCommands(ctx context.Context, output chan<- Command) {
	cmd := UpdateTable{ID{GroupTradingDates, GroupTradingDates}}
	select {
	case output <- &cmd:
	case <-ctx.Done():
	}
}

func prepareLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("не удалось загрузить часовой пояс Москвы")
	}
	return loc
}

var zoneMoscow = prepareLocation()

func nextDayEnd() time.Duration {
	now := time.Now().In(zoneMoscow)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, zoneMoscow)

	if end.After(now) {
		return end.Sub(now)
	}
	return end.AddDate(0, 0, 1).Sub(now)
}

type DayStarted struct {
}

func (d DayStarted) StartProduceCommands(ctx context.Context, output chan<- Command) {
	cmd := UpdateTable{ID{GroupTradingDates, GroupTradingDates}}
	timer := time.NewTimer(nextDayEnd())

LOOP:
	for {

		select {
		case <-timer.C:
			go func() {
				output <- &cmd
			}()
			timer.Reset(nextDayEnd())
		case <-ctx.Done():
			timer.Stop()
			break LOOP
		}
	}
}
