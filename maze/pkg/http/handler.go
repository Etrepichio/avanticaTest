package http

import (
	"context"
	"encoding/json"
	"github.com/avanticaTest/maze/pkg/endpoints"
	"github.com/avanticaTest/maze/pkg/errors"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerFinalizer(RequestLogFinalizer(logger)),
	}

	c := mux.NewRouter()
	c.Methods("POST").Path( "/spot").Handler(httptransport.NewServer(
		endpoints.CreateSpotEndpoint,
		DecodeCreateSpotRequest,
		EncodeCreateSpotResponse,
		append(options)...,
	))
	c.Methods("GET").Path( "/spot/{id}").Handler(httptransport.NewServer(
		endpoints.GetSingleSpotEndpoint,
		DecodeGetSingleSpotRequest,
		EncodeGetSingleSpotResponse,
		append(options)...,
	))
	c.Methods("PUT").Path( "/spot/{id}").Handler(httptransport.NewServer(
		endpoints.ModifySpotEndpoint,
		DecodeModifySpotRequest,
		EncodeModifySpotResponse,
		append(options)...,
	))
	c.Methods("GET").Path( "/spots").Handler(httptransport.NewServer(
		endpoints.GetSpotsEndpoint,
		DecodeGetSpotsRequest,
		EncodeGetSpotsResponse,
		append(options)...,
	))
	c.Methods("DELETE").Path( "/spot/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteSpotEndpoint,
		DecodeDeleteSpotRequest,
		EncodeDeleteSpotResponse,
		append(options)...,
	))

	//PATH endpoints

	c.Methods("POST").Path( "/path").Handler(httptransport.NewServer(
		endpoints.CreatePathEndpoint,
		DecodeCreatePathRequest,
		EncodeCreatePathResponse,
		append(options)...,
	))
	c.Methods("GET").Path( "/path/{id}").Handler(httptransport.NewServer(
		endpoints.GetSinglePathEndpoint,
		DecodeGetSinglePathRequest,
		EncodeGetSinglePathResponse,
		append(options)...,
	))
	c.Methods("PUT").Path( "/path/{id}").Handler(httptransport.NewServer(
		endpoints.ModifyPathEndpoint,
		DecodeModifyPathRequest,
		EncodeModifyPathResponse,
		append(options)...,
	))
	c.Methods("GET").Path( "/paths").Handler(httptransport.NewServer(
		endpoints.GetPathsEndpoint,
		DecodeGetPathsRequest,
		EncodeGetPathsResponse,
		append(options)...,
	))
	c.Methods("DELETE").Path( "/path/{id}").Handler(httptransport.NewServer(
		endpoints.DeletePathEndpoint,
		DecodeDeletePathRequest,
		EncodeDeletePathResponse,
		append(options)...,
	))

	//ORIGIN endpoints

	c.Methods("POST").Path( "/origin").Handler(httptransport.NewServer(
		endpoints.CreateOriginEndpoint,
		DecodeCreateOriginRequest,
		EncodeCreateOriginResponse,
		append(options)...,
	))
	c.Methods("GET").Path( "/origin").Handler(httptransport.NewServer(
		endpoints.GetOriginEndpoint,
		DecodeGetOriginRequest,
		EncodeGetOriginResponse,
		append(options)...,
	))
	c.Methods("PUT").Path( "/origin").Handler(httptransport.NewServer(
		endpoints.ModifyOriginEndpoint,
		DecodeModifyOriginRequest,
		EncodeModifyOriginResponse,
		append(options)...,
	))
	c.Methods("POST").Path( "/quadrantSpots").Handler(httptransport.NewServer(
		endpoints.GetSpotsInQuadrantEndpoint,
		DecodeGetSpotsInQuadrantRequest,
		EncodeGetSpotsInQuadrantResponse,
		append(options)...,
	))
	c.Methods("DELETE").Path( "/origin").Handler(httptransport.NewServer(
		endpoints.DeleteOriginEndpoint,
		DecodeDeleteOriginRequest,
		EncodeDeleteOriginResponse,
		append(options)...,
	))

	return c
}

// RequestLogFinalizer is called at the end of an http request. Use it to log final
// information regarding a request.
func RequestLogFinalizer(logger log.Logger) httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		// log a bunch of response values
		level.Info(logger).Log("code", code,
			"ua", r.UserAgent(),
			"remote", r.RemoteAddr,
			"method", r.Method,
			"url", r.URL,
			"path", r.URL.Path,
			"host", r.Host)
	}
}

// DecodeCreateSpotRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeCreateSpotRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	var rp models.Spot
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.CreateSpotRequest{
		Req: rp,
	}, err
}

// EncodeGetRandomAisleResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeCreateSpotResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.CreateObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeGetSingleSpotRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetSingleSpotRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	pvars := mux.Vars(r)

	id := pvars["id"]

	return endpoints.GetSingleObjectRequest{
		ObjectID: id,
	}, err
}

// EncodeGetSingleSpotResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetSingleSpotResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.GetSingleSpotResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeGetSpotsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetSpotsRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	return endpoints.EmptyGetRequest{}, err
}

// EncodeGetSpotsResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetSpotsResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.GetSpotsResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeModifySpotRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeModifySpotRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	pvars := mux.Vars(r)

	id := pvars["id"]
	var rp models.Spot
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.ModifySpotRequest{
		Req: rp,
		ID:  id,
	}, err
}

// EncodeModifySpotResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeModifySpotResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.ModifyObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeDeleteSpotRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeDeleteSpotRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	pvars := mux.Vars(r)

	id := pvars["id"]

	return endpoints.GetSingleObjectRequest{
		ObjectID: id,
	}, err
}

// EncodeDeleteSpotResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeDeleteSpotResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.ModifyObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}



//Path Decoders / Encoders

// DecodeCreateOriginRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeCreateOriginRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	var rp models.Origin
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.CreateOriginRequest{
		Req: rp,
	}, err
}

// EncodeCreateOriginResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeCreateOriginResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.CreateObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeGetOriginRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetOriginRequest(_ context.Context, r *http.Request) (req interface{}, err error) {


	return endpoints.EmptyGetRequest{}, err
}

// EncodeGetOriginResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetOriginResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.GetOriginResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeGetSpotsInQuadrantRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetSpotsInQuadrantRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	var rp models.Quadrant
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.GetSpotsInQuadrantRequest{
		Req: rp,
	}, err
}

// EncodeGetSpotsInQuadrantResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetSpotsInQuadrantResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.GetSpotsResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeModifyPathRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeModifyOriginRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	var rp models.Origin
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.CreateOriginRequest{
		Req: rp,
	}, err
}

// EncodeModifyOriginResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeModifyOriginResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.ModifyObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeCreateSpotRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeDeleteOriginRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	return endpoints.EmptyGetRequest{}, err
}

// EncodeGetRandomAisleResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeDeleteOriginResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.ModifyObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}


// Origin Endpoints

// DecodeCreatePathRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeCreatePathRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	var rp models.CreatePathRequest
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.CreatePathRequest{
		Req: rp,
	}, err
}

// EncodeCreatePathResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeCreatePathResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.CreateObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeGetSinglePathRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetSinglePathRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	pvars := mux.Vars(r)

	id := pvars["id"]

	return endpoints.GetSingleObjectRequest{
		ObjectID: id,
	}, err
}

// EncodeGetSinglePathResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetSinglePathResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.GetSinglePathResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeGetPathsRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeGetPathsRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	return endpoints.EmptyGetRequest{}, err
}

// EncodeGetPathsResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeGetPathsResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.GetPathsResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeModifyPathRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeModifyPathRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	pvars := mux.Vars(r)

	id := pvars["id"]
	var rp models.CreatePathRequest
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		if err == io.EOF {
			return nil, errors.ErrMissingBodyContent
		} else if err == io.ErrUnexpectedEOF {
			return nil, errors.ErrMalformedBodyContent
		} else {
			return nil, err
		}
	}
	return endpoints.ModifyPathRequest{
		Req: rp,
		ID:  id,
	}, err
}

// EncodeModifyPathResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeModifyPathResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.ModifyObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}

// DecodeCreateSpotRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func DecodeDeletePathRequest(_ context.Context, r *http.Request) (req interface{}, err error) {

	pvars := mux.Vars(r)

	id := pvars["id"]

	return endpoints.GetSingleObjectRequest{
		ObjectID: id,
	}, err
}

// EncodeGetRandomAisleResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeDeletePathResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	// set response header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// cast response to known type
	res, ok := response.(endpoints.ModifyObjectResponse)
	if !ok {
		return errors.ErrResponseEncoding
	}
	if res.Err != nil {
		return res.Err
	}

	// create json
	return json.NewEncoder(w).Encode(res.Res)
}
