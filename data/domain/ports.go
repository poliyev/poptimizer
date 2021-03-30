package domain

import (
	"context"
	"time"
)

type (
	Group string
	Name  string
)

type Table interface {
	Group() Group
	Name() Name
	Update(ctx context.Context, cmd Command) (Event, error)
}

type Event struct {
	Group   Group
	Name    Name
	Replace bool
	Rows    interface{}
}

type Command struct {
	Group   Group
	Name    Name
	LastDay time.Time
}

type Factory interface {
	NewTable(group Group, name Name) Table
}

type Service interface {
	Start(ctx context.Context) <-chan Command
}

type Rule interface {
	Match(event Event) bool
	Handle(event Event) []Command
}
