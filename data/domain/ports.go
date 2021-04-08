package domain

import (
	"context"
)

// Типы для идентификаторов таблиц.
type (
	Group string
	Name  string
)

// TableID используется для идентификации таблиц, команд и событий, связанных с ними.
type TableID struct {
	Group Group
	Name  Name
}

func (i TableID) ID() TableID {
	return i
}

// Identifiable - талицы и связанные с ними объекты.
type Identifiable interface {
	ID() TableID
}

// Event - события, произошедшие при попытке обновить таблицу.
type Event interface {
	Identifiable
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
	NewTable(id TableID) Table
}
