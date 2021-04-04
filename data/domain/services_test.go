package domain

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZone(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")
	assert.Equal(t, prepareZone("Europe/Moscow"), moscow)
}

func TestZonePanic(t *testing.T) {
	assert.Panics(t, func() { prepareZone("WrongZone") })
}

func TestBeforeNextISSDailyUpdate(t *testing.T) {
	in := time.Date(2021, 4, 3, 21, 44, 0, 0, time.UTC)
	out := time.Date(2021, 4, 4, 0, 45, 0, 0, zoneMoscow)

	assert.Equal(t, out, nextISSDailyUpdate(in))
}

func TestAfterNextISSDailyUpdate(t *testing.T) {
	in := time.Date(2021, 4, 3, 21, 46, 0, 0, time.UTC)
	out := time.Date(2021, 4, 5, 0, 45, 0, 0, zoneMoscow)

	assert.Equal(t, out, nextISSDailyUpdate(in))
}

func TestTradingDayAppStart(t *testing.T) {
	out := UpdateTable{ID{groupTradingDates, groupTradingDates}}

	output := make(chan Command)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		CheckTradingDay{}.StartProduceCommands(ctx, output)
	}()

	assert.Equal(t, &out, <-output)

	cancel()
	wg.Wait()
}

func TestTradingDayNextUpdate(t *testing.T) {
	out := UpdateTable{ID{groupTradingDates, groupTradingDates}}

	timer := make(chan time.Time)

	output := make(chan Command)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		CheckTradingDay{timer}.StartProduceCommands(ctx, output)
	}()

	<-output

	// После публикации данных на ISS должна отправляться команда
	now := time.Now()
	timer <- nextISSDailyUpdate(now).Add(time.Second)
	assert.Equal(t, &out, <-output)

	// До начала следующего дня обновление таймера не порождает команд
	close(output)
	timer <- nextISSDailyUpdate(now).Add(time.Hour * 24)

	cancel()
	wg.Wait()
}
