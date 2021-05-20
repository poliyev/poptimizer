package adapters

import (
	"context"
	"errors"
	"fmt"
	"poptimizer/data/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo обеспечивает хранение и загрузку таблиц.
type Repo struct {
	uri    string
	dbName string
	db     *mongo.Database
}

// NewRepo - создает новое Repo.
func NewRepo(mongoURI, mongoDB string) *Repo {
	repo := Repo{
		uri:    mongoURI,
		dbName: mongoDB,
	}

	return &repo
}

// Start запускает модуль репозитория.
func (r *Repo) Start(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.uri))
	if err != nil {
		return fmt.Errorf("fail to start Repo: %w", err)
	}

	r.db = client.Database(r.dbName)

	return nil
}

// Shutdown останавливает модуль репозитория.
func (r *Repo) Shutdown(ctx context.Context) error {
	if err := r.db.Client().Disconnect(ctx); err != nil {
		return fmt.Errorf("repo shutdown error: %w", r.db.Client().Disconnect(ctx))
	}

	return nil
}

// Unmarshal заполняет шаблон таблицы из события.
func (r *Repo) Unmarshal(ctx context.Context, event domain.UpdateRequired) (domain.Table, error) {
	collection := r.db.Collection(string(event.Group()))
	err := collection.FindOne(ctx, bson.M{"_id": event.Name()}).Decode(event.Template)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return event.Template, nil
	case err == nil:
		return event.Template, nil
	default:
		return nil, fmt.Errorf("repo load error %s: %w", event, err)
	}
}

// Replace замещает сохраненные строки таблицы на новые из события.
func (r *Repo) Replace(ctx context.Context, event domain.RowsReplaced) error {
	collection := r.db.Collection(string(event.Group()))

	filter := bson.M{"_id": event.Name()}
	update := bson.M{"$set": bson.M{"rows": event.Rows}}

	if _, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true)); err != nil {
		return fmt.Errorf("repo save error %s: %w", event, err)
	}

	return nil
}

// Append дописывает в конец новые строки из события.
func (r *Repo) Append(ctx context.Context, event domain.RowsAppended) error {
	collection := r.db.Collection(string(event.Group()))

	filter := bson.M{"_id": event.Name()}
	update := bson.M{"$push": bson.M{"rows": bson.M{"$each": event.Rows}}}

	if _, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true)); err != nil {
		return fmt.Errorf("repo save error %s: %w", event, err)
	}

	return nil
}

// ViewJSON загружает ExtendedJSON представление строк из таблицы.
func (r *Repo) ViewJSON(ctx context.Context, id domain.ID) ([]byte, error) {
	collection := r.db.Collection(string(id.Group()))

	projections := options.FindOne().SetProjection(bson.M{"_id": 0, "rows": 1})

	raw, err := collection.FindOne(ctx, bson.M{"_id": id.Name()}, projections).DecodeBytes()
	if err != nil {
		return nil, fmt.Errorf("repo json viewer error %s: %w", id, err)
	}

	json, err := bson.MarshalExtJSON(raw, true, true)
	if err != nil {
		return nil, fmt.Errorf("repo json viewer error %s: %w", id, err)
	}

	return json, nil
}
