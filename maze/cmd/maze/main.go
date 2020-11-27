package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/avanticaTest/maze/pkg/config"
	"github.com/avanticaTest/maze/pkg/db"
	"github.com/avanticaTest/maze/pkg/endpoints"
	"github.com/avanticaTest/maze/pkg/service/path"
	"github.com/avanticaTest/maze/pkg/service/quadrant"
	"github.com/avanticaTest/maze/pkg/service/spot"
	"github.com/go-kit/kit/log"
	"net/http"

	mazehttp "github.com/avanticaTest/maze/pkg/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

func main() {
	var (
		consulAddr = flag.String("consul.addr", "localhost:8500", "Consul agent address")
	)
	flag.Parse()

	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	conf := config.NewConfig("maze_config")
	err := conf.Load(*consulAddr)
	if err != nil {
		panic(err)
	}
	connectionURI := conf.DBURI
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

	spot := spot.New(logger, db.New(client, logger))
	path := path.New(logger, db.New(client, logger))
	quadrant := quadrant.New(logger, db.New(client, logger))

	eps := endpoints.New(spot, path, quadrant, logger)
	handler := mazehttp.NewHTTPHandler(eps, logger)

	http.ListenAndServe(":8080", handler)

}
