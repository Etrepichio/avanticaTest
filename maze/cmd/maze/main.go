package main

import (
	"context"
	"fmt"
	"github.com/avanticaTest/maze/pkg/service/path"
	"github.com/avanticaTest/maze/pkg/service/quadrant"
	"github.com/avanticaTest/maze/pkg/service/spot"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

func main() {

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	connectionURI := "mongodb://maze.db:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	spot := spot.New(logger, client)
	path := path.New(logger, client)
	quadrant := quadrant.New(logger, client)

	router := mux.NewRouter()
	router.HandleFunc("/spot", spot.CreateSpot).Methods("POST")
	router.HandleFunc("/spot/{id}", spot.GetSingleSpot).Methods("GET")
	router.HandleFunc("/spot/{id}", spot.ModifySpot).Methods("PUT")
	router.HandleFunc("/spots", spot.GetSpots).Methods("GET")
	router.HandleFunc("/spot/{id}", spot.DeleteSpot).Methods("DELETE")

	router.HandleFunc("/path", path.CreatePath).Methods("POST")
	router.HandleFunc("/path/{id}", path.GetSinglePath).Methods("GET")
	router.HandleFunc("/path/{id}", path.ModifyPath).Methods("PUT")
	router.HandleFunc("/paths", path.GetPaths).Methods("GET")
	router.HandleFunc("/path/{id}", path.DeletePath).Methods("DELETE")

	router.HandleFunc("/origin", quadrant.CreateOrigin).Methods("POST")
	router.HandleFunc("/origin", quadrant.GetOrigin).Methods("GET")
	router.HandleFunc("/origin", quadrant.ModifyOrigin).Methods("PUT")
	router.HandleFunc("/origin", quadrant.DeleteOrigin).Methods("DELETE")
	router.HandleFunc("/quadrantSpots", quadrant.GetSpotsInQuadrant).Methods("POST")
	http.ListenAndServe(":8080", router)

}
