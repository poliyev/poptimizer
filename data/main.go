package main

import (
	"time"
)

func main() {
	cfg := Config{
		StartTimeout:     time.Minute,
		ShutdownTimeout:  time.Minute,
		RequestTimeout:   time.Microsecond * 600,
		EventBusTimeouts: time.Minute,
		ServerAddr:       "localhost:3000",
		ISSMaxCons:       20,
		MongoURI:         "mongodb://localhost:27017",
		MongoDB:          "new_data",
	}

	app := NewApp(cfg)
	app.Run()
}
