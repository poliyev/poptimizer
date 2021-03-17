package main

import (
	"context"
	"fmt"
	"poptimizer/data/app"
	"poptimizer/data/tables"
)

func main() {
	repo, _ := app.InitRepo()
	table, _ := repo.Get(context.Background(), "trading_dates", "trading_dates")
	cmd := tables.Command{}
	fmt.Println(table.Update(context.Background(), cmd))
	_ = repo.Save(context.Background(), table)
	fmt.Println(table.Update(context.Background(), cmd))
}
