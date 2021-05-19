package domain

import (
	"context"
	"fmt"
)

// Типы для идентификаторов таблиц.
type (
	// Group - группа таблиц. Иногда группа может состоять из одной таблицы.
	Group string
	// Name - название таблицы в рамках группы. Если группа состоит из одной таблицы, то имя должно совпадать названием группы.
	Name string
)

// Identifiable - талицы и связанные с ними объекты.
type Identifiable interface {
	fmt.Stringer
	Group() Group
	Name() Name
}

// ID используется для идентификации таблиц и событий, связанных с ними.
type ID struct {
	group Group
	name  Name
}

// NewID создает идентификатор таблицы.
func NewID(group, name string) ID {
	return ID{Group(group), Name(name)}
}

func (i ID) Group() Group {
	return i.group
}

func (i ID) Name() Name {
	return i.name
}

func (i ID) String() string {
	return fmt.Sprintf("ID(%s, %s)", i.Group(), i.Name())
}

// Event - события, произошедшие при попытке обновить таблицу.
type Event interface {
	Identifiable
}

// Table - таблица, которая умеет обновляться.
type Table interface {
	Identifiable
	Update(ctx context.Context) []Event
}

// Rule - бизнес правило.
//
// Читает события из входящего канала, обрабатывает их с заданных таймаутом и пишет новые события в исходящий канал.
type Rule interface {
	Activate(ctx context.Context, in <-chan Event, out chan<- Event)
}
