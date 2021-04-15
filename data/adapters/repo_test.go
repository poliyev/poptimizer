package adapters

import (
	"context"
	"poptimizer/data/domain"
	"testing"
	"time"

	"github.com/WLM1ke/gomoex"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func prepareRepo(t *testing.T) *Repo {
	t.Helper()

	iss := NewISSClient(20)
	factory := domain.NewMainFactory(iss)

	repo := NewRepo("mongodb://localhost:27017", "test", factory)
	if repo.Start(context.Background()) != nil {
		t.Error("Не удалось запустить тестовую репозиторий.")
	}

	return repo
}

func cleanRepo(t *testing.T, repo *Repo) {
	t.Helper()

	ctx := context.Background()
	if repo.db.Drop(ctx) != nil {
		t.Error("Не удалось удалить тестовую базу.")
	}

	if repo.Shutdown(ctx) != nil {
		t.Error("Не удалось завершить работу репозитория.")
	}
}

var (
	testID  = domain.NewTableID("trading_dates", "trading_dates")
	testRow = gomoex.Date{From: time.Time{}, Till: time.Time{}.AddDate(1, 0, 0)}
)

func TestRepoLoadAbsentTable(t *testing.T) {
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

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
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	event := domain.RowsReplaced{TableID: testID, Rows: []gomoex.Date{testRow}}

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
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	event := domain.RowsAppended{TableID: testID, Rows: []gomoex.Date{testRow}}

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
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	json, err := repo.ViewJSON(context.Background(), testID)

	assert.Nil(t, json)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestRepoJsonWithDoc(t *testing.T) {
	out := `{"rows":[{"from":{"$date":{"$numberLong":"-62135596800000"}},"till":{"$date":{"$numberLong":"-62104060800000"}}}]}`

	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	event := domain.RowsReplaced{TableID: testID, Rows: []gomoex.Date{testRow}}
	if repo.Save(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	json, err := repo.ViewJSON(context.Background(), testID)

	assert.Equal(t, []byte(out), json)
	assert.Nil(t, err)
}
