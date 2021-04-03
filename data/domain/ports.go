package domain

import (
	"context"
)

// Типы для идентификаторов таблиц.
type (
	Group string
	Name  string
)

// Identifiable - талицы и связанные с ними объекты.
type Identifiable interface {
	Group() Group
	Name() Name
}

// Event - события, произошедшие при попытке обновить таблицу.
type Event interface {
	Identifiable
}

// TableUpdated - событие успешного обновления талицы, содержит информацию об измененных строках.
type TableUpdated interface {
	Event
	Rows() interface{}
}

// ErrOccurred - событие ошибки при обновлении таблицы, содержит произошедшую ошибку.
type ErrOccurred interface {
	Event
	Error() error
}

// EventConsumer - обработчик событий.
type EventConsumer interface {
	StartHandleEvent(ctx context.Context, source <-chan Event)
}

// Command - команда для таблицы.
type Command interface {
	Identifiable
}

// CommandSource - генератор команд для таблиц.
type CommandSource interface {
	StartProduceCommands(ctx context.Context, output chan<- Command)
}

// Rule - описывает бизнес правило, принимает произошедшие события и при необходимости генерирует необходимые команды для таблиц.
type Rule interface {
	EventConsumer
	CommandSource
}

// Table - таблица, которая умеет обрабатывать команды и возвращать произошедшие в процессе их исполнения события.
type Table interface {
	Identifiable
	HandleCommand(ctx context.Context, cmd Command) []Event
}

// Factory - фабрика для создания таблиц.
type Factory interface {
	NewTable(group Group, name Name) Table
}
