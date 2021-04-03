package domain

// ID используется для идентификации таблиц, команд и событий, связанных с ними.
type ID struct {
	group Group
	name  Name
}

// Group - группа талицы.
func (id *ID) Group() Group {
	return id.group
}

// Name - имя таблицы в группе.
func (id *ID) Name() Name {
	return id.name
}
