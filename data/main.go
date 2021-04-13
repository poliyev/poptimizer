package main

import (
	cfg2 "poptimizer/data/cfg"
	"time"
)

func main() {
	cfg := cfg2.Config{
		StartTimeout:     time.Minute,
		ShutdownTimeout:  time.Minute,
		RequestTimeout:   time.Microsecond * 600,
		EventBusTimeouts: time.Minute,
		ServerAddr:       "localhost:3000",
		ISSMaxCons:       20,
		MongoURI:         "mongodb://localhost:27017",
		MongoDB:          "new_data",
	}

	app := cfg2.NewApp(cfg)
	app.Run()
}
