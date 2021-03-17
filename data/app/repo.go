package app

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"poptimizer/data/tables"
)

type Repo struct {
	factory tables.Factory
	db      *mongo.Database
}

func (r *Repo) Get(ctx context.Context, group tables.Group, name tables.Name) (tables.Table, error) {
	template := r.factory.NewTable(group, name)
	collection := r.db.Collection(string(group))

	err := collection.FindOne(ctx, bson.M{"_id": name}).Decode(template)
	if err == mongo.ErrNoDocuments {
		return template, nil
	} else if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *Repo) Save(ctx context.Context, table tables.Table) error {
	collection := r.db.Collection(string(table.Group()))

	if _, err := collection.UpdateOne(ctx, bson.M{"_id": table.Name()}, bson.M{"$set": table}, options.Update().SetUpsert(true)); err != nil {
		return err
	}
	return nil
}

func InitRepo() (*Repo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	repo := Repo{
		factory: tables.InitTableFactory(),
		db:      client.Database("new_data"),
	}
	return &repo, nil
}
