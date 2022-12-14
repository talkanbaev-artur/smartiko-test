package server

import (
	"context"
	"ehdw/smartiko-test/src/logic/endpoint"
	"ehdw/smartiko-test/src/logic/service"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	logkit "github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	transp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type errMsg struct {
	Err string `json:"error"`
}

func MakeServer(r *mux.Router, s service.Service) error {
	ends := endpoint.MakeEdnpoints(s)

	options := []transp.ServerOption{
		transp.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(errMsg{err.Error()})
		}),
		transp.ServerErrorHandler(transport.NewLogErrorHandler(logkit.NewLogger(logrus.New()))),
	}

	r.Methods("DELETE").Path("/device/{id}").Handler(
		transp.NewServer(
			ends.DeleteDeviceEndpoint,
			decodeDeviceURL,
			encodeResponse,
			options...,
		),
	)

	r.Methods("DELETE").Path("/devices").Handler(
		transp.NewServer(
			ends.DeleteDevicesEndpoint,
			decodeDeviceArray,
			encodeResponse,
			options...,
		),
	)

	r.Methods("POST").Path("/device").Handler(
		transp.NewServer(
			ends.AddDeviceEndpoint,
			decodeDeviceID,
			encodeCreated,
			options...,
		),
	)

	r.Methods("POST").Path("/devices").Handler(
		transp.NewServer(
			ends.AddDevicesEndpoint,
			decodeDeviceArray,
			encodeCreated,
			options...,
		),
	)

	r.Methods("GET").Path("/device/{id}").Handler(
		transp.NewServer(
			ends.GetDeviceEndpoint,
			decodeDeviceURL,
			encodeDevice,
			options...,
		),
	)

	r.Methods("GET").Path("/devices").Handler(
		transp.NewServer(
			ends.GetDevicesEndpoint,
			decodeEmpty,
			encodeResponse,
			options...,
		),
	)

	return nil
}

func decodeDeviceID(_ context.Context, r *http.Request) (any, error) {
	req := endpoint.DeviceID{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeDeviceArray(_ context.Context, r *http.Request) (any, error) {
	req := []endpoint.DeviceID{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeDeviceURL(_ context.Context, r *http.Request) (any, error) {
	req, ok := endpoint.DeviceID{}, false
	params := mux.Vars(r)
	req.ID, ok = params["id"]
	if !ok {
		return nil, errors.New("id not supplied")
	}
	return req, nil
}

func encodeCreated(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	return json.NewEncoder(w).Encode(response)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeDevice(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if reflect.ValueOf(response).IsZero() {
		w.WriteHeader(404)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
func decodeEmpty(_ context.Context, r *http.Request) (interface{}, error) {
	return "", nil
}
