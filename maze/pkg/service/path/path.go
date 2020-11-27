package path

import (
	"context"
	"github.com/avanticaTest/maze/pkg/db"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func New(logger log.Logger, db db.DBManager) PathHandler {
	return &stubPathHandler{
		db:     db,
		logger: logger,
	}
}

//CreatePath creates a path given two spots ID
func (s *stubPathHandler) CreatePath(ctx context.Context, request models.CreatePathRequest) (string, error) {

	//first we get the objectIDs from the strings of the request
	idpa, _ := primitive.ObjectIDFromHex(request.PointA)
	idpb, _ := primitive.ObjectIDFromHex(request.PointB)
	var path models.Path
	var spotA models.Spot
	var spotB models.Spot
	//then we get the spots represented by those IDs
	err := s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpa}, &spotA)
	if err != nil {
		level.Error(s.logger).Log("method", "CreatePath", "error", err)
		return "", err
	}
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpb}, &spotB)
	if err != nil {
		level.Error(s.logger).Log("method", "CreatePath", "error", err)
		return "", err
	}
	//then we calculate the distance between those points and load it into the path
	path.Distance = Distance(spotA, spotB)
	path.PointA = spotA.ID
	path.PointB = spotB.ID
	//and save the path itself
	result, err := s.db.InsertOne(ctx, "mazedb", "paths", path)
	if err != nil {
		level.Error(s.logger).Log("method", "CreatePath", "error", err)
		return "", err
	}

	return result, nil
}

//GetPaths returns all the paths we have
func (s *stubPathHandler) GetPaths(ctx context.Context) ([]models.Path, error) {

	result, err := s.db.FindPaths(ctx, "mazedb", "paths")
	if err != nil{
		level.Error(s.logger).Log("method", "GetPaths", "error", err)
		return nil, err
	}

	return result, nil
}

//GetSinglePath gets one single path given its ID
func (s *stubPathHandler) GetSinglePath(ctx context.Context, id string)(models.Path, error) {

	var path models.Path

	idp, _ := primitive.ObjectIDFromHex(id)


	err := s.db.FindOne(ctx, "mazedb", "paths", models.Path{ID: idp}, &path)
	if err != nil {
		level.Error(s.logger).Log("method", "GetSinglePath", "error", err)
		return models.Path{},err
	}

	//Here we update the path just in case any of the spots has changed
	var spotA models.Spot
	var spotB models.Spot
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: path.PointA}, &spotA)
	if err != nil {
		level.Error(s.logger).Log("method", "GetSinglePath", "error", err)
		return models.Path{},err
	}
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: path.PointB}, &spotB)
	if err != nil {
		level.Error(s.logger).Log("method", "GetSinglePath", "error", err)
		return models.Path{},err
	}

	path.Distance = Distance(spotA, spotB)
	s.logger.Log("PathA", path.PointA, "PathB", path.PointB)
	return path, nil
}


//ModifyPath modifies a path changing one or both of the spots that compose it
func (s *stubPathHandler) ModifyPath(ctx context.Context, request models.CreatePathRequest, id string) (int, error) {

	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		level.Error(s.logger).Log("method", "ModifyPath", "error", err)
		return 0,err
	}
	filter := bson.D{{"_id", idp}}


	idpa, _ := primitive.ObjectIDFromHex(request.PointA)
	idpb, _ := primitive.ObjectIDFromHex(request.PointB)
	var spotA models.Spot
	var spotB models.Spot
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpa}, &spotA)
	if err != nil {
		level.Error(s.logger).Log("method", "ModifyPath", "error", err)
		return 0,err
	}
	err = s.db.FindOne(ctx, "mazedb", "spots", models.Spot{ID: idpb}, &spotB)
	if err != nil {
		level.Error(s.logger).Log("method", "ModifyPath", "error", err)
		return 0,err
	}


	update := bson.D{{"$set", bson.D{{"point_a", idpa},
		{"point_b", idpb},
		{"distance", Distance(spotA, spotB)}}}}
	result, err := s.db.UpdateOne(ctx, filter, update, "mazedb", "paths")
	if err != nil {
		level.Error(s.logger).Log("method", "GetSpotsInQuadrant", "error", err)
		return 0, err
	}
	return result, nil
}


//DeletePath deletes one path given its ID
func (s *stubPathHandler) DeletePath(ctx context.Context, id string) (int, error) {


	idp, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		level.Error(s.logger).Log("method", "GetSpotsInQuadrant", "error", err)
		return 0, err
	}
	filter := bson.D{{"_id", idp}}

	result, err := s.db.DeleteOne(ctx, filter, "mazedb", "paths")
	if err != nil {
		level.Error(s.logger).Log("method", "DeletePath", "error", err)
		return 0, err
	}

	return result, nil
}

//Distance calculates the distance between two spots
func Distance(a, b models.Spot) float64 {
	first := math.Pow(b.XCoordinate-a.XCoordinate, 2)
	second := math.Pow(b.YCoordinate-a.YCoordinate, 2)
	return math.Sqrt(first + second)

}
