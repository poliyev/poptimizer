package domain

// RowsReplaced - событие замены всех строк в таблице.
type RowsReplaced struct {
	TableID
	Rows interface{}
}

// RowsAppended - событие добавления строк в конец таблицы.
type RowsAppended struct {
	TableID
	Rows interface{}
}
