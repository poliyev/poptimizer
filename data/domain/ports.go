package domain

import (
	"context"
)

type (
	Group string
	Name  string
)

type Identifiable interface {
	Group() Group
	Name() Name
}

type EventProvider interface {
	Events() []Event
}

type Table interface {
	Identifiable
	EventProvider
	Update(ctx context.Context, cmd Command)
}

type Command interface {
	Identifiable
}

type Event interface {
	Identifiable
}

type TableUpdated interface {
	Event
	Rows() interface{}
}

type ErrOccurred interface {
	Identifiable
	Error() error
}

type Factory interface {
	NewTable(group Group, name Name) Table
}

type Service interface {
	Start(ctx context.Context) <-chan Command
}

type Rule interface {
	Match(event Event) bool
	Commands(event Event) []Command
}
