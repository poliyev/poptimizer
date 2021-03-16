package tables

import (
	"context"
	"errors"
	"time"
)

type (
	Group string
	Name  string
	Row   interface{}
)

// Identifiable относится к определенному объекту в группе.
type Identifiable interface {
	Group() Group
	Name() Name
}

// ID базовая реализация
type ID struct {
	group Group
	name  Name
}

func (i *ID) Group() Group {
	return i.group
}

func (i *ID) Name() Name {
	return i.name
}

// Event содержит информации о типе и содержании изменения в новой версии таблицы.
type Event struct {
	ID
	Timestamp   time.Time
	ReplaceRows bool
	NewRows     []Row
}

// Command содержит наименование таблицы, которая должна обработать событие.
type Command interface {
	Identifiable
}

// UpdateTable - основная команда по обновлению таблицы.
type UpdateTable struct {
	ID
}

// CommandHandler обрабатывает команду.
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// Table поддерживает обработку команды и является источником событий.
type Table interface {
	Identifiable
	FlushEvents() []Event
	CommandHandler
}

var ErrRowsValidationErr = errors.New("ошибка валидации данных")

// tableTemplate хранит строки таблицы, реализует этапы обновления данных.
type tableTemplate interface {
	updateCond(timestamp time.Time) bool
	prepareRows(ctx context.Context, cmd Command) error
	validateRows() error
	addNewRows() []Row
	replace() bool
}

// BaseTable хранит метаинформацию о таблице, реализует управление процедурой обновления.
type BaseTable struct {
	ID
	timestamp time.Time
	tableTemplate
	events []Event
}

func (t *BaseTable) Handle(ctx context.Context, cmd Command) (err error) {
	if t.updateCond(t.timestamp) {
		err = t.prepareRows(ctx, cmd)
		if err != nil {
			return err
		}

		err = t.validateRows()
		if err != nil {
			return err
		}

		t.timestamp = time.Now()
		newEvent := Event{
			ID:          t.ID,
			Timestamp:   t.timestamp,
			ReplaceRows: t.replace(),
			NewRows:     t.addNewRows(),
		}
		t.events = append(t.events, newEvent)
	}
	return nil
}

func (t *BaseTable) FlushEvents() []Event {
	events := t.events
	t.events = make([]Event, 0)
	return events
}
