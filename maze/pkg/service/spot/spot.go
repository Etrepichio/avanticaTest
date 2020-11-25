package spot

import (
	"context"
	"github.com/avanticaTest/maze/pkg/internal/db"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpotHandler interface {
	CreateSpot(ctx context.Context, request models.Spot) (string, error)
	GetSingleSpot(ctx context.Context, id string) (models.Spot, error)
	ModifySpot(ctx context.Context, request models.Spot, id string) (int, error)
	GetSpots(ctx context.Context) ([]models.Spot, error)
	DeleteSpot(ctx context.Context, id string) (int, error)
}

type stubSpotHandler struct {
	db     db.DBManager
	logger log.Logger
}

func New(logger log.Logger, client *mongo.Client) SpotHandler {
	return stubSpotHandler{
		db:     db.New(client, logger),
		logger: logger,
	}
}

func (s stubSpotHandler) CreateSpot(ctx context.Context, request models.Spot) (string, error) {

	result, err := s.db.InsertOne(ctx, "mazedb", "spots", request)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s stubSpotHandler) GetSpots(ctx context.Context) ([]models.Spot, error) {

	result, err := s.db.FindSpots(ctx, "mazedb", "spots")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s stubSpotHandler) GetSingleSpot(ctx context.Context, id string) (models.Spot, error) {
	var spot models.Spot
	idp, _ := primitive.ObjectIDFromHex(id)


	err := s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idp}, &spot)
	if err != nil {
		return models.Spot{}, err
	}

	return spot, nil
}

func (s stubSpotHandler) ModifySpot(ctx context.Context, request models.Spot, id string) (int, error) {

	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{"_id", idp}}

	update := bson.D{{"$set", bson.D{{"x_coordinate", request.XCoordinate},
		{"y_coordinate", request.YCoordinate},
		{"name", request.Name},
		{"number", request.Number}}}}
	result, err := s.db.UpdateOne(ctx, filter, update, "mazedb", "spots")
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s stubSpotHandler) DeleteSpot(ctx context.Context, id string) (int, error) {

	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{"_id", idp}}

	result, err := s.db.DeleteOne(ctx, filter, "mazedb", "spots")
	if err != nil {
		return 0, err
	}

	return result, nil
}
