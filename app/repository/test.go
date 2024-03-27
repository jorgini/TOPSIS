package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"webApp/lib/eval"
)

type Test struct {
	ID     uint        `bson:"id"`
	Rating eval.Rating `bson:"rating"`
}

type TestDao struct {
	c *mongo.Collection
}

var TestDb TestDao

func (t *TestDao) ConnectTestCollection(database *mongo.Database) {
	t.c = database.Collection(os.Getenv("TESTCOL"))
}

func (t *TestDao) Insert(ctx context.Context, test *Test) error {
	_, err := t.c.InsertOne(ctx, test)
	return err
}

func (t *TestDao) Get(ctx context.Context, test *Test) (*Test, error) {
	filter := bson.D{{"id", test.ID}}
	var findInst Test
	err := t.c.FindOne(ctx, filter).Decode(&findInst)

	return &findInst, err
}

func (t *TestDao) Delete(ctx context.Context, test *Test) error {
	filter := bson.D{{"id", test.ID}}

	result, err := t.c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
