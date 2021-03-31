package domain

type ID struct {
	group Group
	name  Name
}

func (id *ID) Group() Group {
	return id.group
}

func (id *ID) Name() Name {
	return id.name
}
