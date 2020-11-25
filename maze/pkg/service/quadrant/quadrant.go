package quadrant

import (
	"context"
	"errors"
	"github.com/avanticaTest/maze/pkg/internal/db"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OriginHandler interface {
	CreateOrigin(ctx context.Context, request models.Origin) (string, error)
	ModifyOrigin(ctx context.Context, request models.Origin) (int, error)
	GetOrigin(ctx context.Context) (models.Origin, error)
	GetSpotsInQuadrant(ctx context.Context, request models.Quadrant) ([]models.Spot, error)
	DeleteOrigin(ctx context.Context) (int, error)
}

type stubOriginHandler struct {
	db     db.DBManager
	logger log.Logger
}

func New(logger log.Logger, client *mongo.Client) OriginHandler {
	return stubOriginHandler{
		db:     db.New(client, logger),
		logger: logger,
	}
}

func (s stubOriginHandler) CreateOrigin(ctx context.Context, request models.Origin) (string, error) {

	e, err := s.db.EstimatedDocumentCount(ctx, "mazedb", "origin")
	if err != nil{
		return "", err
	}
	if e > 0{
		return "", errors.New("There can be only one Origin")
	}

	result, err := s.db.InsertOne(ctx, "mazedb", "origin", request)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s stubOriginHandler) GetSpotsInQuadrant(ctx context.Context, request models.Quadrant) ([]models.Spot, error) {

	//first we get all the spots
	resultS, err := s.db.FindSpots(ctx, "mazedb", "spots")
	if err != nil {
		return nil, err
	}



	//then we get the origin

	resultO, err := s.db.FindOrigin(ctx, "mazedb", "origin")
	if err != nil {
		return nil, err
	}
	if len(resultO) < 1{
		return nil, errors.New("Inexistant Origin")
	}
	origin := resultO[0]
	var result []models.Spot
	for _,v := range resultS{

		switch request.Quadrant {
		case "upper_left":
			if (origin.XOrigin >= v.XCoordinate) && (origin.YOrigin <= v.YCoordinate) {
				result = append(result, v)
			}
		case "upper_right":
			if (origin.XOrigin <= v.XCoordinate) && (origin.YOrigin <= v.YCoordinate) {
				result = append(result, v)
			}
		case "bottom_left":
			if (origin.XOrigin >= v.XCoordinate) && (origin.YOrigin >= v.YCoordinate) {
				result = append(result, v)
			}
		case "bottom_right":
			if (origin.XOrigin <= v.XCoordinate) && (origin.YOrigin >= v.YCoordinate) {
				result = append(result, v)
			}
		}

	}
	return result, nil
}

func (s stubOriginHandler) GetOrigin(ctx context.Context) (models.Origin, error) {

	result, err := s.db.FindOrigin(ctx, "mazedb", "origin")
	if err != nil {
		return models.Origin{}, err
	}
	if len(result) > 0{
		return result[0],nil
	}else{
		return models.Origin{}, errors.New("No origin found")
	}

}

func (s stubOriginHandler) ModifyOrigin(ctx context.Context, request models.Origin) (int, error) {


	filter := bson.D{}

	update := bson.D{{"$set", bson.D{{"x_origin", request.XOrigin},
		{"y_origin", request.YOrigin}}}}
	result, err := s.db.UpdateOne(ctx, filter, update, "mazedb", "origin")
	if err != nil {
		return 0, err
	}
	return result, nil
}


func (s stubOriginHandler) DeleteOrigin(ctx context.Context) (int, error) {

	filter := bson.D{}

	result, err := s.db.DeleteMany(ctx, filter, "mazedb", "origin")
	if err != nil {
		return 0, err
	}

	return result, nil
}
