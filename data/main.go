package main

import (
	"poptimizer/data/config"
	"time"
)

const (
	issMaxCons     = 20
	serverTimeouts = time.Millisecond * 40
)

func main() {
	cfg := &config.Config{
		StartTimeout:     time.Minute,
		ShutdownTimeout:  time.Minute,
		ServerAddr:       "localhost:3000",
		ServerTimeouts:   serverTimeouts,
		EventBusTimeouts: time.Minute,
		ISSMaxCons:       issMaxCons,
		MongoURI:         "mongodb://localhost:27017",
		MongoDB:          "new_data",
	}

	app := config.NewApp(cfg)
	app.Run()
}
