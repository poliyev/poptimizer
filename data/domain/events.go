package domain

// RowsReplaced - событие замены всех строк в таблице.
type RowsReplaced struct {
	id
	rows interface{}
}

// Rows - новые строки таблицы.
func (r *RowsReplaced) Rows() interface{} {
	return r.rows
}

// RowsAppended - событие добавления строк в конец таблицы.
type RowsAppended struct {
	id
	rows interface{}
}

// Rows - добавленные строки в таблицу.
func (r *RowsAppended) Rows() interface{} {
	return r.rows
}
