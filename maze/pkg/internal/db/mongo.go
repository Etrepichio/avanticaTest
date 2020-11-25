package db

import (
	"context"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBManager interface {
	FindOne(ctx context.Context, db, collection string, filter, objective interface{}) error
	InsertOne(ctx context.Context, db, collection string, doc interface{}) (string, error)
	UpdateOne(ctx context.Context, filter, update interface{}, db, collection string) (int, error)
	DeleteOne(ctx context.Context, filter interface{}, db, collection string) (int, error)
	DeleteMany(ctx context.Context, filter interface{}, db, col string) (int, error)
	FindSpots(ctx context.Context, db, col string) (result []models.Spot, err error)
	FindPaths(ctx context.Context, db, col string) (result []models.Path, err error)
	FindOrigin(ctx context.Context, db, col string) (result []models.Origin, err error)
	EstimatedDocumentCount(ctx context.Context, db, collection string) (int, error)
}

type stubDBManager struct {
	client *mongo.Client
	logger log.Logger
}

func New(client *mongo.Client, logger log.Logger) DBManager {
	return &stubDBManager{
		client: client,
		logger: logger,
	}
}

func (s *stubDBManager) FindOne(ctx context.Context, db, col string, filter, objective interface{}) error {

	collection := s.client.Database(db).Collection(col)

	err := collection.FindOne(ctx, filter).Decode(objective)
	if err != nil {
		return err
	}
	return nil
}

func (s *stubDBManager) InsertOne(ctx context.Context, db, col string, doc interface{}) (string, error) {

	collection := s.client.Database(db).Collection(col)

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	oid := res.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

func (s *stubDBManager) FindSpots(ctx context.Context, db, col string) (result []models.Spot, err error) {

	collection := s.client.Database(db).Collection(col)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Spot
		cursor.Decode(&p)
		result = append(result, p)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *stubDBManager) FindPaths(ctx context.Context, db, col string) (result []models.Path, err error) {

	collection := s.client.Database(db).Collection(col)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Path
		cursor.Decode(&p)
		result = append(result, p)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *stubDBManager) FindOrigin(ctx context.Context, db, col string) (result []models.Origin, err error) {

	collection := s.client.Database(db).Collection(col)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Origin
		cursor.Decode(&p)
		result = append(result, p)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *stubDBManager) UpdateOne(ctx context.Context, filter, update interface{}, db, col string) (int, error) {

	collection := s.client.Database(db).Collection(col)

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return int(result.ModifiedCount), nil
}

func (s *stubDBManager) DeleteOne(ctx context.Context, filter interface{}, db, col string) (int, error) {

	collection := s.client.Database(db).Collection(col)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}


func (s *stubDBManager) DeleteMany(ctx context.Context, filter interface{}, db, col string) (int, error) {

	collection := s.client.Database(db).Collection(col)

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}


func (s *stubDBManager) EstimatedDocumentCount(ctx context.Context, db, col string) (int, error){

	collection := s.client.Database(db).Collection(col)

	r, err := collection.EstimatedDocumentCount(ctx)
	if err != nil{
		return 0, err
	}
	return int(r),nil
}