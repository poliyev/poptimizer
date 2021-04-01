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
	StartHandleEvent(ctx context.Context, source <-chan Event)
}

type Command interface {
	Identifiable
}

type CommandSource interface {
	StartProduceCommands(ctx context.Context, output chan<- Command)
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
