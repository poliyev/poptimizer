package adapters

import (
	"context"
	"github.com/WLM1ke/gomoex"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"poptimizer/data/domain"
	"testing"
	"time"
)

func prepareRepo() *repo {
	iss := NewISSClient(20)
	factory := domain.NewMainFactory(iss)
	repo := NewRepo("mongodb://localhost:27017", "test", factory)
	repo.Start(context.Background())
	return repo
}

func cleanRepo(repo *repo) {
	ctx := context.Background()
	repo.db.Drop(ctx)
	repo.Shutdown(ctx)
}

var testID = domain.TableID{"trading_dates", "trading_dates"}
var testRow = gomoex.Date{time.Time{}, time.Time{}.AddDate(1, 0, 0)}

func TestRepoLoadAbsentTable(t *testing.T) {
	repo := prepareRepo()
	defer cleanRepo(repo)

	table, err := repo.Load(context.Background(), testID)
	if err != nil {
		t.Error("Не удалось загрузить таблицу")
		return
	}

	dates, ok := table.(*domain.TradingDates)
	if !ok {
		t.Error("Некорректная таблица")
		return
	}
	assert.Nil(t, dates.Rows)
}

func TestRepoSaveReplaceEvent(t *testing.T) {
	repo := prepareRepo()
	defer cleanRepo(repo)

	event := domain.RowsReplaced{testID, []gomoex.Date{testRow}}

	if repo.Save(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}
	if repo.Save(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	table, err := repo.Load(context.Background(), testID)
	if err != nil {
		t.Error("Не удалось загрузить сохраненную таблицу")
		return
	}

	dates, ok := table.(*domain.TradingDates)
	if !ok {
		t.Error("Некорректная таблица")
		return
	}

	rows := dates.Rows
	assert.Equal(t, 1, len(rows))
	assert.Equal(t, testRow, rows[0])

}

func TestRepoSaveAppendEvent(t *testing.T) {
	repo := prepareRepo()
	defer cleanRepo(repo)

	event := domain.RowsAppended{testID, []gomoex.Date{testRow}}

	if repo.Save(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}
	if repo.Save(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}
	table, err := repo.Load(context.Background(), testID)
	if err != nil {
		t.Error("Не удалось загрузить сохраненную таблицу")
		return
	}

	dates, ok := table.(*domain.TradingDates)
	if !ok {
		t.Error("Некорректная таблица")
		return
	}

	rows := dates.Rows
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, testRow, rows[0])

}

func TestRepoJsonNoDoc(t *testing.T) {
	repo := prepareRepo()
	defer cleanRepo(repo)

	json, err := repo.ViewJSON(context.Background(), testID)

	assert.Nil(t, json)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestRepoJsonWithDoc(t *testing.T) {
	out := `{"rows":[{"from":{"$date":{"$numberLong":"-62135596800000"}},"till":{"$date":{"$numberLong":"-62104060800000"}}}]}`

	repo := prepareRepo()
	defer cleanRepo(repo)

	event := domain.RowsReplaced{testID, []gomoex.Date{testRow}}
	if repo.Save(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	json, err := repo.ViewJSON(context.Background(), testID)

	assert.Equal(t, []byte(out), json)
	assert.Nil(t, err)
}
