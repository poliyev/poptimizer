package domain

// id используется для идентификации таблиц, команд и событий, связанных с ними.
type id struct {
	group Group
	name  Name
}

// Group - группа талицы.
func (id *id) Group() Group {
	return id.group
}

// Name - имя таблицы в группе.
func (id *id) Name() Name {
	return id.name
}
