package main

import (
	"context"
	"fmt"
	"github.com/avanticaTest/pkg/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)



func main() {

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodbhost"))
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	maze := service.New(logger, client)
	defer client.Disconnect(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/spot", maze.CreateSpot).Methods("POST")
	http.ListenAndServe(":8080", router)

	//database := client.Database("avantica")
	//spotsCollection := database.Collection("spots")
	//pathCollection := database.Collection("paths")
	//
	//fmt.Println(spotsCollection)
	//fmt.Println(pathCollection)

}
