package path

import (
	"context"
	"github.com/avanticaTest/maze/pkg/internal/db"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

type PathHandler interface {
	CreatePath(ctx context.Context, request models.CreatePathRequest) (string, error)
	ModifyPath(ctx context.Context, request models.CreatePathRequest, id string) (int, error)
	GetSinglePath(ctx context.Context, id string)(models.Path, error)
	GetPaths(ctx context.Context) ([]models.Path, error)
	DeletePath(ctx context.Context, id string) (int, error)
}

type stubPathHandler struct {
	db     db.DBManager
	logger log.Logger
}

func New(logger log.Logger, client *mongo.Client) PathHandler {
	return &stubPathHandler{
		db:     db.New(client, logger),
		logger: logger,
	}
}


func (s *stubPathHandler) CreatePath(ctx context.Context, request models.CreatePathRequest) (string, error) {

	idpa, _ := primitive.ObjectIDFromHex(request.PointA)
	idpb, _ := primitive.ObjectIDFromHex(request.PointB)
	var path models.Path
	var spotA models.Spot
	var spotB models.Spot
	err := s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpa}, &spotA)
	if err != nil {
		return "", err
	}
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpb}, &spotB)
	if err != nil {
		return "", err
	}
	path.Distance = Distance(spotA, spotB)
	path.PointA = spotA.ID
	path.PointB = spotB.ID
	result, err := s.db.InsertOne(ctx, "mazedb", "paths", path)
	if err != nil {
		return "", err
	}

	return result, nil
}


func (s *stubPathHandler) GetPaths(ctx context.Context) ([]models.Path, error) {

	result, err := s.db.FindPaths(ctx, "mazedb", "paths")
	if err != nil{
		return nil, err
	}

	return result, nil
}


func (s *stubPathHandler) GetSinglePath(ctx context.Context, id string)(models.Path, error) {

	var path models.Path

	idp, _ := primitive.ObjectIDFromHex(id)


	err := s.db.FindOne(ctx, "mazedb", "paths", models.Path{ID: idp}, &path)
	if err != nil {
		return models.Path{},err
	}

	//Here we update the path just in case any of the spots has changed
	var spotA models.Spot
	var spotB models.Spot
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: path.PointA}, &spotA)
	if err != nil {
		return models.Path{},err
	}
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: path.PointB}, &spotB)
	if err != nil {
		return models.Path{},err
	}

	path.Distance = Distance(spotA, spotB)
	s.logger.Log("PathA", path.PointA, "PathB", path.PointB)
	return path, nil
}



func (s *stubPathHandler) ModifyPath(ctx context.Context, request models.CreatePathRequest, id string) (int, error) {


	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return 0,err
	}
	filter := bson.D{{"_id", idp}}


	idpa, _ := primitive.ObjectIDFromHex(request.PointA)
	idpb, _ := primitive.ObjectIDFromHex(request.PointB)
	var spotA models.Spot
	var spotB models.Spot
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpa}, &spotA)
	if err != nil {
		return 0,err
	}
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpb}, &spotB)
	if err != nil {
		return 0,err
	}


	update := bson.D{{"$set", bson.D{{"point_a", idpa},
		{"point_b", idpb},
		{"distance", Distance(spotA, spotB)}}}}
	result, err := s.db.UpdateOne(ctx, filter, update, "mazedb", "paths")
	if err != nil {
		return 0, err
	}
	return result, nil
}



func (s *stubPathHandler) DeletePath(ctx context.Context, id string) (int, error) {


	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return 0, err
	}
	filter := bson.D{{"_id", idp}}

	result, err := s.db.DeleteOne(ctx, filter, "mazedb", "paths")
	if err != nil {
		return 0, err
	}

	return result, nil
}

func Distance(a, b models.Spot) float64 {
	first := math.Pow(b.XCoordinate-a.XCoordinate, 2)
	second := math.Pow(b.YCoordinate-a.YCoordinate, 2)
	return math.Sqrt(first + second)

}
