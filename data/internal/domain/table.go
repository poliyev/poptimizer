package domain

// Table представляет таблицу с данными, актуальными на конкретную дату.
type Table[R any] struct {
	Version
	Rows []R
}

func (t Table[R]) IsEmpty() bool {
	return len(t.Rows) == 0
}

func (t Table[R]) LastRow() R {
	return t.Rows[len(t.Rows)-1]
}
