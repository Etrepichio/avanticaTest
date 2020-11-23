package service

import (
	"context"
	"encoding/json"
	"github.com/avanticaTest/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type MazeService interface {
	CreateSpot(w http.ResponseWriter, r *http.Request)
}

type stubMazeHandler struct {
	db     *mongo.Client
	logger *logrus.Logger
}

func New(logger *logrus.Logger, db *mongo.Client) MazeService {
	return stubMazeHandler{
		db:     db,
		logger: logger,
	}
}

func (s stubMazeHandler) CreateSpot(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var spot models.Spot
	if err := json.NewDecoder(request.Body).Decode(&spot); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	collection := s.db.Database("avantica").Collection("spots")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, spot)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}
