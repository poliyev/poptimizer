package config

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestModule struct {
	name string
	step chan<- string
}

func (t TestModule) Name() string {
	return t.name
}

func (t TestModule) Start(ctx context.Context) error {
	t.step <- fmt.Sprintf("%s start", t.name)

	return nil
}

func (t TestModule) Shutdown(ctx context.Context) error {
	t.step <- fmt.Sprintf("%s shutdown", t.name)

	return nil
}

func TestAppSteps(t *testing.T) {
	step := make(chan string, 4)
	moduleA := TestModule{"A", step}
	moduleB := TestModule{"B", step}

	stop := make(chan os.Signal, 1)
	app := App{time.Second, time.Second, stop, []Module{moduleA, moduleB}}

	go app.Run()

	stop <- syscall.SIGINT

	assert.Equal(t, "A start", <-step)
	assert.Equal(t, "B start", <-step)
	assert.Equal(t, "B shutdown", <-step)
	assert.Equal(t, "A shutdown", <-step)
}
