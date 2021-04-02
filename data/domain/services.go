package domain

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

func nextDayStart() time.Time {
	now := time.Now().In(zoneMoscow)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 45, 0, 0, zoneMoscow)

	if end.Before(now) {
		end.AddDate(0, 0, 1).Sub(now)
	}
	return end
}

type DayStarted struct {
}

func (d DayStarted) StartProduceCommands(ctx context.Context, output chan<- Command) {
	cmd := UpdateTable{ID{GroupTradingDates, GroupTradingDates}}

	output <- &cmd
	nextDay := nextDayStart()

	timer := time.Tick(time.Hour)
	wake := make(chan os.Signal, 1)
	signal.Notify(wake, syscall.SIGCONT)

LOOP:
	for {
		select {
		case <-timer:
			if time.Now().After(nextDay) {
				output <- &cmd
				nextDay = nextDayStart()
			}
		case <-wake:
			fmt.Printf("Программа проснулась после сна")
			if time.Now().After(nextDay) {
				output <- &cmd
				nextDay = nextDayStart()
			}
		case <-ctx.Done():
			break LOOP
		}
	}
}
