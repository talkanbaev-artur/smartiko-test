package endpoint

import (
	"context"
	"ehdw/smartiko-test/src/logic/service"
	"ehdw/smartiko-test/src/util"

	logkit "github.com/go-kit/kit/log/logrus"
	"github.com/sirupsen/logrus"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddDeviceEndpoint     endpoint.Endpoint
	AddDevicesEndpoint    endpoint.Endpoint
	GetDeviceEndpoint     endpoint.Endpoint
	GetDevicesEndpoint    endpoint.Endpoint
	DeleteDeviceEndpoint  endpoint.Endpoint
	DeleteDevicesEndpoint endpoint.Endpoint
}

func MakeEdnpoints(s service.Service) Endpoints {
	es := Endpoints{}
	es.AddDeviceEndpoint = makeAddDeviceEndpoint(s)
	es.AddDevicesEndpoint = makeAddDevicesEndpoint(s)
	es.DeleteDeviceEndpoint = makeDeleteDeviceEndpoint(s)
	es.GetDeviceEndpoint = makeGetDeviceEndpoint(s)
	es.GetDevicesEndpoint = makeGetDevicesEndpoint(s)
	es.DeleteDevicesEndpoint = makeDeleteDevicesEndpoint(s)
	return es
}

type DeviceID struct {
	ID string `json:"id"`
}

func makeAddDeviceEndpoint(s service.Service) endpoint.Endpoint {
	logger := logkit.NewLogger(logrus.New())
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeviceID)
		response, err = s.AddDevice(ctx, req.ID)
		return
	}
	e = util.LoggingMiddleware(logger, "add device")(e)
	return e
}

func makeAddDevicesEndpoint(s service.Service) endpoint.Endpoint {
	logger := logkit.NewLogger(logrus.New())
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]DeviceID)
		resp := []int{}
		for _, r := range req {
			var id int
			id, err = s.AddDevice(ctx, r.ID)
			if err != nil {
				return
			}
			resp = append(resp, id)
		}
		response = resp
		return
	}
	e = util.LoggingMiddleware(logger, "add devices")(e)
	return e
}

func makeGetDeviceEndpoint(s service.Service) endpoint.Endpoint {
	logger := logkit.NewLogger(logrus.New())
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeviceID)
		return s.GetDevice(ctx, req.ID)
	}
	e = util.LoggingMiddleware(logger, "get device")(e)
	return e
}

func makeGetDevicesEndpoint(s service.Service) endpoint.Endpoint {
	logger := logkit.NewLogger(logrus.New())
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetDevices(ctx)
	}
	e = util.LoggingMiddleware(logger, "get devices")(e)
	return e
}

func makeDeleteDeviceEndpoint(s service.Service) endpoint.Endpoint {
	logger := logkit.NewLogger(logrus.New())
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeviceID)
		err = s.DeleteDevice(ctx, req.ID)
		return
	}
	e = util.LoggingMiddleware(logger, "delete device")(e)
	return e
}

func makeDeleteDevicesEndpoint(s service.Service) endpoint.Endpoint {
	logger := logkit.NewLogger(logrus.New())
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]DeviceID)
		for _, r := range req {
			err = s.DeleteDevice(ctx, r.ID)
			if err != nil {
				return
			}
		}
		return
	}
	e = util.LoggingMiddleware(logger, "delete devices")(e)
	return e
}
