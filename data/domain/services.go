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

func lastDay() time.Time {
	now := time.Now().In(zoneMoscow)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, zoneMoscow)

	days := 1
	if end.After(now) {
		days = 2
	}
	end = end.AddDate(0, 0, -days)

	return time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
}

type WorkStarted struct {
	out chan Command
}

func (d *WorkStarted) Start(ctx context.Context) {
	d.out = make(chan Command)

	go func() {
		cmd := UpdateTable{ID{GroupTradingDates, GroupTradingDates}}
		select {
		case d.out <- &cmd:
			close(d.out)
		case <-ctx.Done():
		}
	}()

}

func (d *WorkStarted) Commands() <-chan Command {
	return d.out
}
