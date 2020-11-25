package endpoints

import (
	"context"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/avanticaTest/maze/pkg/service/path"
	"github.com/avanticaTest/maze/pkg/service/quadrant"
	"github.com/avanticaTest/maze/pkg/service/spot"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CreateSpotEndpoint    endpoint.Endpoint
	GetSingleSpotEndpoint endpoint.Endpoint
	GetSpotsEndpoint      endpoint.Endpoint
	ModifySpotEndpoint    endpoint.Endpoint
	DeleteSpotEndpoint    endpoint.Endpoint

	CreatePathEndpoint    endpoint.Endpoint
	GetSinglePathEndpoint endpoint.Endpoint
	GetPathsEndpoint      endpoint.Endpoint
	ModifyPathEndpoint    endpoint.Endpoint
	DeletePathEndpoint    endpoint.Endpoint

	CreateOriginEndpoint       endpoint.Endpoint
	GetOriginEndpoint          endpoint.Endpoint
	GetSpotsInQuadrantEndpoint endpoint.Endpoint
	ModifyOriginEndpoint       endpoint.Endpoint
	DeleteOriginEndpoint       endpoint.Endpoint
}

// New will create an Endpoints struct with initialized endpoint(s) and
// middleware(s).
//This allow us to wrap our service's final functions with layers of logging, decoding, encoding, etc, abstracting the core functionality of the service
//from the rest
func New(spot spot.SpotHandler, path path.PathHandler, orig quadrant.OriginHandler, logger log.Logger) (ep Endpoints) {
	// create the GetMinesweeper endpoint

	//Spot Endpoints:

	ep.CreateSpotEndpoint = MakeCreateSpotEndpoint(spot)
	ep.CreateSpotEndpoint = LoggingMiddleware(log.With(logger, "method", "CreateSpot"))(ep.CreateSpotEndpoint)

	//create the NewGame endpoint
	ep.GetSingleSpotEndpoint = MakeGetSingleSpotEndpoint(spot)
	ep.GetSingleSpotEndpoint = LoggingMiddleware(log.With(logger, "method", "GetSingleSpot"))(ep.GetSingleSpotEndpoint)

	//create the LoadGame endpoint
	ep.GetSpotsEndpoint = MakeGetSpotsEndpoint(spot)
	ep.GetSpotsEndpoint = LoggingMiddleware(log.With(logger, "method", "GetSpots"))(ep.GetSpotsEndpoint)

	//create the SaveGame endpoint
	ep.ModifySpotEndpoint = MakeModifySpotEndpoint(spot)
	ep.ModifySpotEndpoint = LoggingMiddleware(log.With(logger, "method", "ModifySpot"))(ep.ModifySpotEndpoint)

	//create the Click endpoint
	ep.DeleteSpotEndpoint = MakeDeleteSpotEndpoint(spot)
	ep.DeleteSpotEndpoint = LoggingMiddleware(log.With(logger, "method", "DeleteSpot"))(ep.DeleteSpotEndpoint)

	//Path Endpoints:

	ep.CreatePathEndpoint = MakeCreatePathEndpoint(path)
	ep.CreatePathEndpoint = LoggingMiddleware(log.With(logger, "method", "CreatePath"))(ep.CreatePathEndpoint)

	//create the NewGame endpoint
	ep.GetSinglePathEndpoint = MakeGetSinglePathEndpoint(path)
	ep.GetSinglePathEndpoint = LoggingMiddleware(log.With(logger, "method", "GetSinglePath"))(ep.GetSinglePathEndpoint)

	//create the LoadGame endpoint
	ep.GetPathsEndpoint = MakeGetPathsEndpoint(path)
	ep.GetPathsEndpoint = LoggingMiddleware(log.With(logger, "method", "GetPaths"))(ep.GetPathsEndpoint)

	//create the SaveGame endpoint
	ep.ModifyPathEndpoint = MakeModifyPathEndpoint(path)
	ep.ModifyPathEndpoint = LoggingMiddleware(log.With(logger, "method", "ModifyPath"))(ep.ModifyPathEndpoint)

	//create the Click endpoint
	ep.DeletePathEndpoint = MakeDeletePathEndpoint(path)
	ep.DeletePathEndpoint = LoggingMiddleware(log.With(logger, "method", "DeletePath"))(ep.DeletePathEndpoint)

	//Origin Endpoints:

	ep.CreateOriginEndpoint = MakeCreateOriginEndpoint(orig)
	ep.CreateOriginEndpoint = LoggingMiddleware(log.With(logger, "method", "CreateOrigin"))(ep.CreateOriginEndpoint)

	//create the NewGame endpoint
	ep.GetOriginEndpoint = MakeGetOriginEndpoint(orig)
	ep.GetOriginEndpoint = LoggingMiddleware(log.With(logger, "method", "GetOrigin"))(ep.GetOriginEndpoint)

	//create the LoadGame endpoint
	ep.GetSpotsInQuadrantEndpoint = MakeGetSpotsInQuadrantEndpoint(orig)
	ep.GetSpotsInQuadrantEndpoint = LoggingMiddleware(log.With(logger, "method", "GetSpotsInQuadrant"))(ep.GetSpotsInQuadrantEndpoint)

	//create the SaveGame endpoint
	ep.ModifyOriginEndpoint = MakeModifyOriginEndpoint(orig)
	ep.ModifyOriginEndpoint = LoggingMiddleware(log.With(logger, "method", "ModifyOrigin"))(ep.ModifyOriginEndpoint)

	//create the Click endpoint
	ep.DeleteOriginEndpoint = MakeDeleteOriginEndpoint(orig)
	ep.DeleteOriginEndpoint = LoggingMiddleware(log.With(logger, "method", "DeleteOrigin"))(ep.DeleteOriginEndpoint)

	return ep

}

// MakeCreateSpotEndpoint returns an endpoint that invokes CreateSpot on the service.
func MakeCreateSpotEndpoint(svc spot.SpotHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateSpotRequest)
		res, err := svc.CreateSpot(ctx, req.Req)

		// wrap service response with endpoint response
		return CreateObjectResponse{Res: models.CreateObjectResponse{ID: res}, Err: err}, nil
	}
}

// MakeGetSingleSpotEndpoint returns an endpoint that invokes GetSingleSpot on the service.
func MakeGetSingleSpotEndpoint(svc spot.SpotHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetSingleObjectRequest)
		res, err := svc.GetSingleSpot(ctx, req.ObjectID)

		// wrap service response with endpoint response
		return GetSingleSpotResponse{Res: res, Err: err}, nil
	}
}

// MakeGetSpotsEndpoint returns an endpoint that invokes GetSpots on the service.
func MakeGetSpotsEndpoint(svc spot.SpotHandler) (ep endpoint.Endpoint) {

	// interface parameter is ignored because request does not
	// require input parameters.
	return func(ctx context.Context, _ interface{}) (interface{}, error) {

		res, err := svc.GetSpots(ctx)

		// wrap service response with endpoint response
		return GetSpotsResponse{Res: res, Err: err}, nil
	}
}

// MakeModifySpotEndpoint returns an endpoint that invokes ModifySpot on the service.
func MakeModifySpotEndpoint(svc spot.SpotHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(ModifySpotRequest)
		res, err := svc.ModifySpot(ctx, req.Req, req.ID)

		// wrap service response with endpoint response
		return ModifyObjectResponse{Res: models.ModifyObjectResponse{AffectedItems: res}, Err: err}, nil
	}
}

// MakeDeleteSpotEndpoint returns an endpoint that invokes DeleteSpot on the service.
func MakeDeleteSpotEndpoint(svc spot.SpotHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetSingleObjectRequest)
		res, err := svc.DeleteSpot(ctx, req.ObjectID)

		// wrap service response with endpoint response
		return ModifyObjectResponse{Res: models.ModifyObjectResponse{AffectedItems: res}, Err: err}, nil
	}
}

//Make Paths Endpoints

// MakeCreatePathEndpoint returns an endpoint that invokes CreatePath on the service.
func MakeCreatePathEndpoint(svc path.PathHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreatePathRequest)
		res, err := svc.CreatePath(ctx, req.Req)

		// wrap service response with endpoint response
		return CreateObjectResponse{Res: models.CreateObjectResponse{ID: res}, Err: err}, nil
	}
}

// MakeGetSinglePathEndpoint returns an endpoint that invokes GetSinglePath on the service.
func MakeGetSinglePathEndpoint(svc path.PathHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetSingleObjectRequest)
		res, err := svc.GetSinglePath(ctx, req.ObjectID)

		// wrap service response with endpoint response
		return GetSinglePathResponse{Res: res, Err: err}, nil
	}
}

// MakeGetPathsEndpoint returns an endpoint that invokes GetPaths on the service.
func MakeGetPathsEndpoint(svc path.PathHandler) (ep endpoint.Endpoint) {

	// interface parameter is ignored because request does not
	// require input parameters.
	return func(ctx context.Context, _ interface{}) (interface{}, error) {

		res, err := svc.GetPaths(ctx)

		// wrap service response with endpoint response
		return GetPathsResponse{Res: res, Err: err}, nil
	}
}

// MakeModifyPathEndpoint returns an endpoint that invokes ModifyPath on the service.
func MakeModifyPathEndpoint(svc path.PathHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(ModifyPathRequest)
		res, err := svc.ModifyPath(ctx, req.Req, req.ID)

		// wrap service response with endpoint response
		return ModifyObjectResponse{Res: models.ModifyObjectResponse{AffectedItems: res}, Err: err}, nil
	}
}

// MakeDeletePathEndpoint returns an endpoint that invokes DeletePath on the service.
func MakeDeletePathEndpoint(svc path.PathHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetSingleObjectRequest)
		res, err := svc.DeletePath(ctx, req.ObjectID)

		// wrap service response with endpoint response
		return ModifyObjectResponse{Res: models.ModifyObjectResponse{AffectedItems: res}, Err: err}, nil
	}
}

//Make Origin Endpoints

// MakeCreateOriginEndpoint returns an endpoint that invokes CreateOrigin on the service.
func MakeCreateOriginEndpoint(svc quadrant.OriginHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateOriginRequest)
		res, err := svc.CreateOrigin(ctx, req.Req)

		// wrap service response with endpoint response
		return CreateObjectResponse{Res: models.CreateObjectResponse{ID: res}, Err: err}, nil
	}
}

// MakeGetOriginEndpoint returns an endpoint that invokes GetOrigin on the service.
func MakeGetOriginEndpoint(svc quadrant.OriginHandler) (ep endpoint.Endpoint) {

	// interface parameter is ignored because request does not
	// require input parameters.
	return func(ctx context.Context, _ interface{}) (interface{}, error) {

		res, err := svc.GetOrigin(ctx)

		// wrap service response with endpoint response
		return GetOriginResponse{Res: res, Err: err}, nil
	}
}

// MakeGetSpotsInQuadrantEndpoint returns an endpoint that invokes GetSpotsInQuadrant on the service.
func MakeGetSpotsInQuadrantEndpoint(svc quadrant.OriginHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetSpotsInQuadrantRequest)
		res, err := svc.GetSpotsInQuadrant(ctx, req.Req)

		// wrap service response with endpoint response
		return GetSpotsResponse{Res: res, Err: err}, nil
	}
}

// MakeModifyOriginEndpoint returns an endpoint that invokes ModifyOrigin on the service.
func MakeModifyOriginEndpoint(svc quadrant.OriginHandler) (ep endpoint.Endpoint) {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateOriginRequest)
		res, err := svc.ModifyOrigin(ctx, req.Req)

		// wrap service response with endpoint response
		return ModifyObjectResponse{Res: models.ModifyObjectResponse{AffectedItems: res}, Err: err}, nil
	}
}

// MakeDeletePathEndpoint returns an endpoint that invokes DeleteOrigin on the service.
func MakeDeleteOriginEndpoint(svc quadrant.OriginHandler) (ep endpoint.Endpoint) {

	// interface parameter is ignored because request does not
	// require input parameters.
	return func(ctx context.Context, _ interface{}) (interface{}, error) {

		res, err := svc.DeleteOrigin(ctx)

		// wrap service response with endpoint response
		return ModifyObjectResponse{Res: models.ModifyObjectResponse{AffectedItems: res}, Err: err}, nil
	}
}

type CreateSpotRequest struct {
	Req models.Spot
}

type CreatePathRequest struct {
	Req models.CreatePathRequest
}

type CreateOriginRequest struct {
	Req models.Origin
}

type CreateObjectResponse struct {
	Res models.CreateObjectResponse
	Err error
}

type GetSingleObjectRequest struct {
	ObjectID string
}

type GetSingleSpotResponse struct {
	Res models.Spot
	Err error
}

type GetSinglePathResponse struct {
	Res models.Path
	Err error
}

type GetOriginResponse struct {
	Res models.Origin
	Err error
}

type GetSpotsResponse struct {
	Res []models.Spot
	Err error
}

type GetPathsResponse struct {
	Res []models.Path
	Err error
}

type GetSpotsInQuadrantRequest struct {
	Req models.Quadrant
}

type ModifySpotRequest struct {
	Req models.Spot
	ID  string
}

type ModifyPathRequest struct {
	Req models.CreatePathRequest
	ID  string
}

type ModifyObjectResponse struct {
	Res models.ModifyObjectResponse
	Err error
}

type EmptyGetRequest struct{}
