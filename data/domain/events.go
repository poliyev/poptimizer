package domain

// UpdateRequired требуется обновление таблицы.
//
// Событие содержит указатель на незаполненную таблицу с правильным ID.
// Событие должно перехватываться на уровне приложения, а данные таблицы загружаться из хранилища.
type UpdateRequired struct {
	Template Table
}

// Group - группа таблицы.
func (u UpdateRequired) Group() Group {
	return u.Template.Group()
}

// Name - имя таблицы в рамках группы.
func (u UpdateRequired) Name() Name {
	return u.Template.Name()
}

// String - текстовое представление ID таблицы.
func (u UpdateRequired) String() string {
	return u.Template.String()
}

// RowsReplaced - событие замены всех строк в таблице.
//
// Событие должно перехватываться на уровне приложения, а измененные данные сохраняться в хранилище.
type RowsReplaced struct {
	ID
	Rows interface{}
}

// RowsAppended - событие добавления строк в конец таблицы.
//
// Событие должно перехватываться на уровне приложения, а измененные данные сохраняться в хранилище.
type RowsAppended struct {
	ID
	Rows interface{}
}
