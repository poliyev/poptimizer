package tables

import (
	"context"
	"errors"
	"time"
)

type (
	TableGroup  string
	TableName   string
	Row         interface{}
	CommandType string
)

// TableID указывает на группу и название таблицы.
type TableID interface {
	Group() TableGroup
	Name() TableName
}

// TableVersion содержит информацию о времени изменения таблицы.
type TableVersion interface {
	TableID
	Timestamp() time.Time
}

// Event содержит информации о типе и содержании изменения в новой версии таблицы.
type Event struct {
	Group       TableGroup
	Name        TableName
	Timestamp   time.Time
	ReplaceRows bool
	NewRows     []Row
}

// EventsSource источник событий выдает новые события, очищая внутреннее хранилище событий.
type EventsSource interface {
	FlushEvents() []Event
}

var ErrWrongCommandType = errors.New("неверный тип команды")
var ErrRowsValidationErr = errors.New("ошибка валидации данных")

// Command содержит наименование таблицы, которая должна обработать событие.
type Command interface {
	TableID
	Type() CommandType
}

// CommandHandler обрабатывает команду.
type CommandHandler interface {
	Handle(context.Context, Command) error
}

// Table поддерживает обработку команды и является источником событий.
type Table interface {
	TableVersion
	EventsSource
	CommandHandler
}
