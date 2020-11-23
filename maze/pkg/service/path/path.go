package path

import (
	"context"
	"encoding/json"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"net/http"
	"time"
)

type PathHandler interface {
	CreatePath(w http.ResponseWriter, r *http.Request)
	ModifyPath(w http.ResponseWriter, r *http.Request)
	GetSinglePath(w http.ResponseWriter, r *http.Request)
	GetPaths(w http.ResponseWriter, r *http.Request)
	DeletePath(w http.ResponseWriter, r *http.Request)
}

type stubPathHandler struct {
	db     *mongo.Client
	logger *logrus.Logger
}

func New(logger *logrus.Logger, db *mongo.Client) PathHandler {
	return stubPathHandler{
		db:     db,
		logger: logger,
	}
}

func (s stubPathHandler) CreatePath(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var path models.Path
	if err := json.NewDecoder(request.Body).Decode(&path); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	var spotA models.Spot
	var spotB models.Spot
	spotCollection := s.db.Database("mazedb").Collection("spots")
	pathCollection := s.db.Database("mazedb").Collection("paths")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := spotCollection.FindOne(ctx, models.Spot{ID: path.PointA}).Decode(&spotA)
	err = spotCollection.FindOne(ctx, models.Spot{ID: path.PointB}).Decode(&spotB)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	path.Distance = Distance(spotA, spotB)

	result, err := pathCollection.InsertOne(ctx, path)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func (s stubPathHandler) GetPaths(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var paths []models.Path
	collection := s.db.Database("mazedb").Collection("paths")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var p models.Path
		cursor.Decode(&p)
		paths = append(paths, p)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(paths)
}

func (s stubPathHandler) GetSinglePath(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var path models.Path
	spotCollection := s.db.Database("mazedb").Collection("spots")
	pathCollection := s.db.Database("mazedb").Collection("paths")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := pathCollection.FindOne(ctx, models.Path{ID: id}).Decode(&path)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	//Here we update the path just in case any of the spots has changed
	var spotA models.Spot
	var spotB models.Spot
	err = spotCollection.FindOne(ctx, models.Spot{ID: path.PointA}).Decode(&spotA)
	err = spotCollection.FindOne(ctx, models.Spot{ID: path.PointB}).Decode(&spotB)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	path.Distance = Distance(spotA, spotB)

	json.NewEncoder(response).Encode(path)
}

func (s stubPathHandler) ModifyPath(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var path models.Path
	if err := json.NewDecoder(request.Body).Decode(&path); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.D{{"_id", id}}

	spotCollection := s.db.Database("mazedb").Collection("spots")
	pathCollection := s.db.Database("mazedb").Collection("paths")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var spotA models.Spot
	var spotB models.Spot
	err := spotCollection.FindOne(ctx, models.Spot{ID: path.PointA}).Decode(&spotA)
	err = spotCollection.FindOne(ctx, models.Spot{ID: path.PointB}).Decode(&spotB)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	update := bson.D{{"$set", bson.D{{"point_a", path.PointA},
		{"point_b", path.PointB},
		{"distance", Distance(spotA, spotB)}}}}
	result, err := pathCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func (s stubPathHandler) DeletePath(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var path models.Path
	if err := json.NewDecoder(request.Body).Decode(&path); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.D{{"_id", id}}

	pathCollection := s.db.Database("mazedb").Collection("paths")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	result, err := pathCollection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func Distance(a, b models.Spot) float64 {
	first := math.Pow(b.XCoordinate-a.XCoordinate, 2)
	second := math.Pow(b.YCoordinate-a.YCoordinate, 2)
	return math.Sqrt(first + second)

}
