package app

import (
	"fmt"
	"poptimizer/data/services"
	"poptimizer/data/tables"
)

type App struct {
}

func (a App) Run() {
	events := make(chan tables.Event)
	services.InitServices(events)

	for event := range events {
		fmt.Println(event)
	}

	//repo, _ := InitRepo()
	//c := time.Tick(time.Minute)
	//fmt.Println("App Started!!!")
	//for t := range c {
	//	table, _ := repo.Get(context.Background(), "trading_dates", "trading_dates")
	//	cmd := tables.Command{}
	//	fmt.Println(t)
	//	fmt.Println(table.Update(context.Background(), cmd))
	//	_ = repo.Save(context.Background(), table)
	//}

}
