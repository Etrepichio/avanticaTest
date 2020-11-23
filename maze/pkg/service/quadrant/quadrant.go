package quadrant

import (
	"context"
	"encoding/json"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type OriginHandler interface {
	CreateOrigin(w http.ResponseWriter, r *http.Request)
	ModifyOrigin(w http.ResponseWriter, r *http.Request)
	GetOrigin(w http.ResponseWriter, r *http.Request)
	GetSpotsInQuadrant(w http.ResponseWriter, r *http.Request)
	DeleteOrigin(w http.ResponseWriter, r *http.Request)
}

type stubOriginHandler struct {
	db     *mongo.Client
	logger *logrus.Logger
}

func New(logger *logrus.Logger, db *mongo.Client) OriginHandler {
	return stubOriginHandler{
		db:     db,
		logger: logger,
	}
}

func (s stubOriginHandler) CreateOrigin(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var orig models.Origin
	if err := json.NewDecoder(request.Body).Decode(&orig); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	collection := s.db.Database("mazedb").Collection("origin")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	a, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if a > 0 {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "message": "There can be only one Origin" }`))
		return
	}

	result, err := collection.InsertOne(ctx, orig)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func (s stubOriginHandler) GetSpotsInQuadrant(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var quadrant models.Quadrant
	if err := json.NewDecoder(request.Body).Decode(&quadrant); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	//first we get all the spots
	var spots []models.Spot
	collection := s.db.Database("mazedb").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)

	//then we get the origin
	var orig models.Origin
	origCollection := s.db.Database("mazedb").Collection("origin")
	ocursor, err := origCollection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	defer ocursor.Close(ctx)
	for ocursor.Next(ctx) {
		ocursor.Decode(&orig)
	}
	for cursor.Next(ctx) {
		var s models.Spot
		cursor.Decode(&s)
		switch quadrant.Quadrant {
		case "upper_left":
			if (orig.XOrigin >= s.XCoordinate) && (orig.YOrigin <= s.YCoordinate) {
				spots = append(spots, s)
			}
		case "upper_right":
			if (orig.XOrigin <= s.XCoordinate) && (orig.YOrigin <= s.YCoordinate) {
				spots = append(spots, s)
			}
		case "bottom_left":
			if (orig.XOrigin >= s.XCoordinate) && (orig.YOrigin >= s.YCoordinate) {
				spots = append(spots, s)
			}
		case "bottom_right":
			if (orig.XOrigin <= s.XCoordinate) && (orig.YOrigin >= s.YCoordinate) {
				spots = append(spots, s)
			}
		}

	}
	json.NewEncoder(response).Encode(spots)
}

func (s stubOriginHandler) GetOrigin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var orig models.Origin
	collection := s.db.Database("mazedb").Collection("origin")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&orig)
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(orig)
}

func (s stubOriginHandler) ModifyOrigin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var orig models.Origin
	if err := json.NewDecoder(request.Body).Decode(&orig); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	filter := bson.D{}

	collection := s.db.Database("mazedb").Collection("origin")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	update := bson.D{{"$set", bson.D{{"x_origin", orig.XOrigin},
		{"y_origin", orig.YOrigin}}}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}


func (s stubOriginHandler) DeleteOrigin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var orig models.Origin
	if err := json.NewDecoder(request.Body).Decode(&orig); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	filter := bson.D{}

	collection := s.db.Database("mazedb").Collection("origin")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
