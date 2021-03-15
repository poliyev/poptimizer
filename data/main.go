package main

import (
	"context"
	"fmt"
	"github.com/WLM1ke/gomoex"
	"net/http"
	"poptimizer/data/tables"
)

func main() {
	iss := gomoex.NewISSClient(http.DefaultClient)
	table := tables.NewTradingDates(iss)
	fmt.Println(table.Timestamp())
	cmd := tables.CheckTradingDates{}
	_ = table.Handle(context.Background(), cmd)
	_ = table.Handle(context.Background(), cmd)
	fmt.Println(table.Timestamp())
	fmt.Println(len(table.FlushEvents()))
}
