package services

import "poptimizer/data/tables"

type Start struct {
	Out chan<- tables.Event
}

func (s Start) Run() {
	go func() { s.Out <- tables.Event{} }()
}
