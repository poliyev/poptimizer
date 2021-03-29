package services

import "poptimizer/data/tables"

func InitServices(events chan<- tables.Event) {
	Start{events}.Run()
}
