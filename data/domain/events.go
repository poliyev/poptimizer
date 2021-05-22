package domain

// UpdateRequired требуется обновление таблицы.
//
// Событие содержит указатели на незаполненные таблицы с правильным ID, которые нужны для построения агрегата, начиная с таблицы,
// являющейся его корнем. Событие должно перехватываться на уровне приложения, а данные таблицы загружаться из хранилища.
type UpdateRequired struct {
	Templates []Table
}

// Group - группа таблицы корня агрегата.
func (u UpdateRequired) Group() Group {
	return u.Templates[0].Group()
}

// Name - имя таблицы корня агрегата в рамках группы.
func (u UpdateRequired) Name() Name {
	return u.Templates[0].Name()
}

// UpdateError - ошибка при обновлении таблицы.
type UpdateError struct {
	ID
	Error error
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
