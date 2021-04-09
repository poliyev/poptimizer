package adapters

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"poptimizer/data/domain"
)

// Repo обеспечивает хранение и загрузку таблиц.
type Repo struct {
	factory domain.Factory
	db      *mongo.Database
}

// Load загружает или возвращает пустую новую таблицу.
func (r *Repo) Load(ctx context.Context, id domain.TableID) (domain.Table, error) {
	template := r.factory.NewTable(id)
	collection := r.db.Collection(string(id.Group))

	err := collection.FindOne(ctx, bson.M{"_id": id.Name}).Decode(template)
	switch {
	case err == mongo.ErrNoDocuments:
		return template, nil
	case err == nil:
		return template, nil
	default:
		return nil, err
	}
}

// ViewJOSN загружает JSON представление строк из таблицы.
func (r *Repo) ViewJOSN(ctx context.Context, id domain.TableID) ([]byte, error) {
	collection := r.db.Collection(string(id.Group))

	projections := options.FindOne().SetProjection(bson.M{"_id": 0, "rows": 1})
	raw, err := collection.FindOne(ctx, bson.M{"_id": id.Name}, projections).DecodeBytes()
	if err != nil {
		return nil, err
	}
	return bson.MarshalExtJSON(raw, true, true)
}

// Save сохраняет результаты изменения таблицы.
func (r *Repo) Save(ctx context.Context, event domain.Event) error {
	fmt.Println(event)
	id := event.ID()
	collection := r.db.Collection(string(id.Group))

	filter := bson.M{"_id": id.Name}

	var update bson.M

	switch changes := event.(type) {
	case domain.RowsReplaced:
		update = bson.M{"$set": bson.M{"rows": changes.Rows}}
	case domain.RowsAppended:
		update = bson.M{"$push": bson.M{"rows": bson.M{"$each": changes.Rows}}}
	}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// NewRepo - создает новое Repo.
func NewRepo(ctx context.Context, mongoURI string, mongoDB string, factory domain.Factory) *Repo {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		zap.L().Panic("Не удалось запустить MongoDB", zap.Error(err))
	}
	zap.L().Info("MongoDB работает", zap.Error(err))

	repo := Repo{
		factory: factory,
		db:      client.Database(mongoDB),
	}

	go func() {
		<-ctx.Done()
		if client.Disconnect(context.Background()) != nil {
			zap.L().Error("Не удалось остановить MongoDB", zap.Error(err))
		} else {
			zap.L().Info("Завершена остановка MongoDB", zap.Error(err))
		}
	}()

	return &repo
}
