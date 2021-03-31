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

type Event interface {
	Identifiable
}

type TableUpdated interface {
	Event
	Rows() interface{}
}

type ErrOccurred interface {
	Event
	Error() error
}

type EventConsumer interface {
	HandleEvent(ctx context.Context, cmd Event)
}

type Command interface {
	Identifiable
}

type CommandSource interface {
	Commands() <-chan Command
}

type Rule interface {
	EventConsumer
	CommandSource
}

type Table interface {
	Identifiable
	HandleCommand(ctx context.Context, cmd Command) []Event
}

type Factory interface {
	NewTable(group Group, name Name) Table
}

type Service interface {
	Start(ctx context.Context)
}
