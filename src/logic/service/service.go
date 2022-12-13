package service

import (
	"context"
	"ehdw/smartiko-test/src/logic/service/model"
	"errors"
	"fmt"
)

type Service interface {
	AddDevice(ctx context.Context, deviceName string) (int, error)
	DeleteDevice(ctx context.Context, deviceName string) error
	GetDevice(ctx context.Context, deviceName string) (model.Device, error)
	GetDevices(ctx context.Context) ([]model.Device, error)
}

type Repository interface {
	CreateDevice(ctx context.Context, devID string) (int, error)
	GetDevice(ctx context.Context, devID string) (model.Device, error)
	GetDevices(ctx context.Context) ([]model.Device, error)
	DeleteDevice(ctx context.Context, devID string) error
	ModifyFlags(ctx context.Context, devID string, flags []model.Flag) error
}

type MessageProcessingFunc func(topicName, msgBody string)

type Queue interface {
	RegisterDevices(ctx context.Context, f MessageProcessingFunc, deviceNames ...string) error
}

type service struct {
	r Repository
	q Queue
}

func NewService(rp Repository, q Queue) Service {
	s := service{r: rp, q: q}
	return s
}

var ErrNotImplemented = errors.New("method not implemented")

type ServiceError struct {
	Err    error
	Stage  string
	Method string
}

func (e ServiceError) Update(stage string, err error) ServiceError {
	e.Err = err
	e.Stage = stage
	return e
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("failed processing of '%s' at the stage '%s'; err: %s", e.Method, e.Stage, e.Err.Error())
}

func (s service) AddDevice(ctx context.Context, deviceName string) (int, error) {
	e := ServiceError{Method: "add device"}
	id, err := s.r.CreateDevice(ctx, deviceName)
	if err != nil {
		return 0, e.Update("db query", err)
	}
	return id, nil
}

func (s service) DeleteDevice(ctx context.Context, deviceName string) error {
	e := ServiceError{Method: "delete device"}
	err := s.r.DeleteDevice(ctx, deviceName)
	if err != nil {
		return e.Update("db query", err)
	}
	return nil
}

func (s service) GetDevice(ctx context.Context, deviceName string) (model.Device, error) {
	e := ServiceError{Method: "get device"}
	dev, err := s.r.GetDevice(ctx, deviceName)
	if err != nil {
		return model.Device{}, e.Update("db query", err)
	}
	return dev, nil
}

func (s service) GetDevices(ctx context.Context) ([]model.Device, error) {
	e := ServiceError{Method: "get devices"}
	devs, err := s.r.GetDevices(ctx)
	if err != nil {
		return nil, e.Update("db query", err)
	}
	return devs, nil
}
