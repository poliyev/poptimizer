package services

import (
	"context"
	"poptimizer/data/tables"
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

type DayStarted struct {
	day time.Time
}

func (d DayStarted) Day() time.Time {
	return d.day
}

func newDayStarted() DayStarted {
	now := time.Now().In(zoneMoscow)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, zoneMoscow)

	days := 1
	if end.After(now) {
		days = 2
	}
	end = end.AddDate(0, 0, -days)

	return DayStarted{time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)}
}

type DayBegins struct {
	Out chan<- tables.Event
}

func (s DayBegins) Run(ctx context.Context) {
	go func() {
		timer := time.Tick(time.Second * 5)
		for {
			select {
			case <-timer:
				s.Out <- newDayStarted()
			case <-ctx.Done():
				return
			}
		}
	}()
}

type DayStartedHandler struct {
	Out chan<- tables.Command
}

func (d DayStartedHandler) Match(event tables.Event) bool {
	switch event.(type) {
	case DayStarted:
		return true
	}
	return false
}

func (d DayStartedHandler) Handle(event tables.Event) {
	d.Out <- tables.UpdateTable{"trading_dates", "trading_dates", event.Day()}
}
