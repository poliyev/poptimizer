package tables

type Factory interface {
	NewTable(group Group, name Name) Table
}

type GroupFactory interface {
	group() Group
	singleton() bool
	NewTable(name Name) Table
}

// TableFactory регистрирует и создает таблицы.
type TableFactory struct {
	groupTemplates map[Group]GroupFactory
}

func (t *TableFactory) registerTemplate(factory GroupFactory) {
	if t.groupTemplates == nil {
		t.groupTemplates = make(map[Group]GroupFactory, 1)
	}
	t.groupTemplates[factory.group()] = factory
}

func (t *TableFactory) NewTable(group Group, name Name) Table {
	factory, ok := t.groupTemplates[group]
	if !ok {
		panic("незарегестрированая группа")
	}

	if factory.singleton() && factory.group() != group {
		panic("некорректное имя таблицы")
	}

	return factory.NewTable(name)
}
