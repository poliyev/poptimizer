package main

import (
	"poptimizer/data/config"
	"time"
)

func main() {
	issMaxCons := 20
	cfg := &config.Config{
		StartTimeout:     time.Minute,
		ShutdownTimeout:  time.Minute,
		ServerAddr:       "localhost:3000",
		ServerTimeouts:   time.Millisecond * 40,
		EventBusTimeouts: time.Minute,
		ISSMaxCons:       issMaxCons,
		MongoURI:         "mongodb://localhost:27017",
		MongoDB:          "new_data",
	}

	app := config.NewApp(cfg)
	app.Run()
}
