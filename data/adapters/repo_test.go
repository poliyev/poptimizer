package adapters

import (
	"context"
	"poptimizer/data/domain"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/WLM1ke/gomoex"
	"github.com/stretchr/testify/assert"
)

func prepareRepo(t *testing.T) *Repo {
	t.Helper()

	repo := NewRepo("mongodb://localhost:27017", "test")
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
	testID  = domain.NewID(domain.GroupTradingDates, domain.GroupTradingDates)
	testRow = gomoex.Date{From: time.Time{}, Till: time.Time{}.AddDate(1, 0, 0)}
)

func TestRepoUnmarshalAbsentTable(t *testing.T) {
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	template := domain.TradingDates{ID: testID}

	err := repo.Unmarshal(context.Background(), &template)
	if err != nil {
		t.Error("Не удалось загрузить таблицу")

		return
	}

	assert.Nil(t, template.Rows)
}

func TestRepoSaveReplaceEvent(t *testing.T) {
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	event := domain.RowsReplaced{ID: testID, Rows: []gomoex.Date{testRow}}

	if repo.Replace(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	if repo.Replace(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	template := domain.TradingDates{ID: testID}

	err := repo.Unmarshal(context.Background(), &template)
	if err != nil {
		t.Error("Не удалось загрузить сохраненную таблицу")

		return
	}

	rows := template.Rows
	assert.Equal(t, 1, len(rows))
	assert.Equal(t, testRow, rows[0])
}

func TestRepoSaveAppendEvent(t *testing.T) {
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	event := domain.RowsAppended{ID: testID, Rows: []gomoex.Date{testRow}}

	if repo.Append(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	if repo.Append(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	template := domain.TradingDates{ID: testID}

	err := repo.Unmarshal(context.Background(), &template)
	if err != nil {
		t.Error("Не удалось загрузить сохраненную таблицу")

		return
	}

	rows := template.Rows
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, testRow, rows[0])
}

func TestRepoJsonNoDoc(t *testing.T) {
	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	json, err := repo.ViewJSON(context.Background(), testID)

	assert.Nil(t, json)
	assert.ErrorIs(t, err, mongo.ErrNoDocuments)
}

func TestRepoJsonWithDoc(t *testing.T) {
	out := `{"rows":[{"from":{"$date":{"$numberLong":"-62135596800000"}},"till":{"$date":{"$numberLong":"-62104060800000"}}}]}`

	repo := prepareRepo(t)
	defer cleanRepo(t, repo)

	event := domain.RowsReplaced{ID: testID, Rows: []gomoex.Date{testRow}}
	if repo.Replace(context.Background(), event) != nil {
		t.Error("Не удалось сохранить таблицу")
	}

	json, err := repo.ViewJSON(context.Background(), testID)

	assert.Equal(t, []byte(out), json)
	assert.Nil(t, err)
}
