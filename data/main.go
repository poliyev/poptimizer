package main

import (
	"poptimizer/data/config"
	"time"
)

func main() {
	cfg := config.Config{
		StartTimeout:     time.Minute,
		ShutdownTimeout:  time.Minute,
		ServerAddr:       "localhost:3000",
		ServerTimeouts:   time.Microsecond * 600,
		EventBusTimeouts: time.Minute,
		ISSMaxCons:       20,
		MongoURI:         "mongodb://localhost:27017",
		MongoDB:          "new_data",
	}

	app := config.NewApp(cfg)
	app.Run()
}
