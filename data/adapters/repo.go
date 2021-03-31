package adapters

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"poptimizer/data/domain"
)

type Repo struct {
	factory domain.Factory
	db      *mongo.Database
}

func (r *Repo) Load(ctx context.Context, group domain.Group, name domain.Name) (domain.Table, error) {
	template := r.factory.NewTable(group, name)
	collection := r.db.Collection(string(group))

	err := collection.FindOne(ctx, bson.M{"_id": name}).Decode(template)

	switch {
	case err == mongo.ErrNoDocuments, err == nil:
		return template, nil
	default:
		return nil, err
	}
}

func (r *Repo) Save(ctx context.Context, event domain.TableUpdated) error {
	collection := r.db.Collection(string(event.Group()))

	filter := bson.M{"_id": event.Name()}

	var update bson.M
	switch event.(type) {
	case *domain.RowsReplaced:
		update = bson.M{"$set": bson.M{"rows": event.Rows()}}
	case *domain.RowsAppended:
		update = bson.M{"$push": bson.M{"rows": bson.M{"$each": event.Rows()}}}
	}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func NewRepo(factory domain.Factory) (*Repo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	repo := Repo{
		factory: factory,
		db:      client.Database("new_data"),
	}
	return &repo, nil
}
