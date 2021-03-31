package domain

type RowsReplaced struct {
	ID
	rows interface{}
}

func (r *RowsReplaced) Rows() interface{} {
	return r.rows
}

type RowsAppended struct {
	ID
	rows interface{}
}

func (r *RowsAppended) Rows() interface{} {
	return r.rows
}

type TableUpdateErrOccurred struct {
	ID
	err error
}

func (r *TableUpdateErrOccurred) Error() error {
	return r.err
}
