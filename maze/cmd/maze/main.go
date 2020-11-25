package main

import (
	"context"
	"fmt"
	"github.com/avanticaTest/maze/pkg/endpoints"
	"github.com/avanticaTest/maze/pkg/service/path"
	"github.com/avanticaTest/maze/pkg/service/quadrant"
	"github.com/avanticaTest/maze/pkg/service/spot"
	"github.com/go-kit/kit/log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	mazehttp "github.com/avanticaTest/maze/pkg/http"
	"os"
	"time"
)

func main() {

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	//prepare MongoDB client
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
	// Prepare logging
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		logger = log.With(logger, "svc", "maze")
	}


	spot := spot.New(logger, client)
	path := path.New(logger, client)
	quadrant := quadrant.New(logger, client)

	eps := endpoints.New(spot, path, quadrant, logger)
	handler := mazehttp.NewHTTPHandler(eps, logger)


	http.ListenAndServe(":8080", handler)

}
