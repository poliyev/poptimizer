package tables

type Factory interface {
	NewTable(group Group, name Name) Table
}

// MainFactory регистрирует и создает таблицы.
type MainFactory struct {
	groupTemplates map[Group]Factory
	singletons     map[Group]bool
}

func (t *MainFactory) registerTemplate(group Group, factory Factory, singleton bool) {
	if t.groupTemplates == nil {
		t.groupTemplates = make(map[Group]Factory, 0)
		t.singletons = make(map[Group]bool, 0)
	}
	t.groupTemplates[group] = factory
	t.singletons[group] = singleton
}

func (t *MainFactory) NewTable(group Group, name Name) Table {
	factory, ok := t.groupTemplates[group]
	if !ok {
		panic("незарегестрированая группа")
	}

	if t.singletons[group] && group != Group(name) {
		panic("некорректное имя таблицы")
	}

	return factory.NewTable(group, name)
}
