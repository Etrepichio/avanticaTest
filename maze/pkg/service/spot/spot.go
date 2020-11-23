package spot

import (
	"context"
	"encoding/json"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type SpotHandler interface {
	CreateSpot(w http.ResponseWriter, r *http.Request)
	GetSingleSpot(w http.ResponseWriter, r *http.Request)
	ModifySpot(w http.ResponseWriter, r *http.Request)
	GetSpots(w http.ResponseWriter, r *http.Request)
	DeleteSpot(w http.ResponseWriter, r *http.Request)
}

type stubSpotHandler struct {
	db     *mongo.Client
	logger *logrus.Logger
}

func New(logger *logrus.Logger, db *mongo.Client) SpotHandler {
	return stubSpotHandler{
		db:     db,
		logger: logger,
	}
}

func (s stubSpotHandler) CreateSpot(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var spot models.Spot
	if err := json.NewDecoder(request.Body).Decode(&spot); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	collection := s.db.Database("mazedb").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, spot)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func (s stubSpotHandler) GetSpots(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var spots []models.Spot
	collection := s.db.Database("mazedb").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var s models.Spot
		cursor.Decode(&s)
		spots = append(spots, s)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(spots)
}

func (s stubSpotHandler) GetSingleSpot(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var spot models.Spot
	collection := s.db.Database("mazedb").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, models.Spot{ID: id}).Decode(&spot)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(spot)
}

func (s stubSpotHandler) ModifySpot(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var spot models.Spot
	if err := json.NewDecoder(request.Body).Decode(&spot); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.D{{"_id", id}}

	collection := s.db.Database("mazedb").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	update := bson.D{{"$set", bson.D{{"x_coordinate", spot.XCoordinate},
		{"y_coordinate", spot.YCoordinate},
		{"name", spot.Name},
		{"number", spot.Number}}}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}


func (s stubSpotHandler) DeleteSpot(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")

	var spot models.Spot
	if err := json.NewDecoder(request.Body).Decode(&spot); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.D{{"_id", id}}

	collection := s.db.Database("mazedb").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}