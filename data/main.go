package main

import (
	"poptimizer/data/app"
	"time"
)

func main() {
	cfg := app.Config{
		StartTimeout:    time.Minute,
		ShutdownTimeout: time.Minute,
		RequestTimeout:  time.Millisecond * 100,
		ServerAddr:      ":3000",
		MongoURI:        "mongodb://localhost:27017",
		MongoDB:         "new_data",
	}

	server := app.NewServer(cfg)
	server.Run()
}
