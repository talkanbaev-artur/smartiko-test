package service

import (
	"context"
	"ehdw/smartiko-test/src/logic/service/model"
	"errors"
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

type service struct {
	r Repository
}

func NewService(rp Repository) Service {
	s := service{r: rp}
	//TODO:
	//start mqtt worker here
	return s
}

var ErrNotImplemented = errors.New("method not implemented")

func (s service) AddDevice(ctx context.Context, deviceName string) (int, error) {
	return 0, ErrNotImplemented
}

func (s service) DeleteDevice(ctx context.Context, deviceName string) error {
	return ErrNotImplemented
}

func (s service) GetDevice(ctx context.Context, deviceName string) (model.Device, error) {
	return model.Device{}, ErrNotImplemented
}

func (s service) GetDevices(ctx context.Context) ([]model.Device, error) {
	return nil, ErrNotImplemented
}
