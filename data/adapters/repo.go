package adapters

import (
	"fmt"
	"poptimizer/data/domain"
)

type Repo struct {
	Factory domain.Factory
}

func (r *Repo) Load(group domain.Group, name domain.Name) domain.Table {
	fmt.Println("Загрузка из репо")
	return r.Factory.NewTable(group, name)
}

func (r *Repo) Save(event domain.Event) {
	fmt.Println("Сохранение в репо")
}
