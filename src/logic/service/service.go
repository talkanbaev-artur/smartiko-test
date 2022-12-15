package service

import (
	"context"
	"ehdw/smartiko-test/src/logic/service/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service interface {
	AddDevice(ctx context.Context, deviceName string) (int, error)
	DeleteDevice(ctx context.Context, deviceName string) error
	GetDevice(ctx context.Context, deviceName string) (*model.Device, error)
	GetDevices(ctx context.Context) ([]model.Device, error)
}

type Repository interface {
	CreateDevice(ctx context.Context, devID string) (int, error)
	GetDevice(ctx context.Context, devID string) (model.Device, error)
	GetDevices(ctx context.Context) ([]model.Device, error)
	DeleteDevice(ctx context.Context, devID string) error
	ModifyFlags(ctx context.Context, devID string, flags []*model.Flag) error
	GetAllEnabledDevices(ctx context.Context) ([]string, error)
}

type MessageProcessingFunc func(topicName, msgBody string)

type Queue interface {
	RegisterDevices(ctx context.Context, f MessageProcessingFunc, deviceNames ...string) error
	UnsubscribeDevices(ctx context.Context, deviceNames ...string) error
}

type service struct {
	r Repository
	q Queue
}

func NewService(rp Repository, q Queue) Service {
	s := service{r: rp, q: q}
	devices, _ := rp.GetAllEnabledDevices(context.Background())
	q.RegisterDevices(context.Background(), s.generalMessageProcessor, devices...)
	return s
}

var ErrNotImplemented = errors.New("method not implemented")
var NoDeviceFound = errors.New("requested device does not exist")

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
	s.q.RegisterDevices(ctx, s.generalMessageProcessor, deviceName)
	return id, nil
}

type flagMsg struct {
	DeviceID string `json:"dev_eui"`
	ParamID  int    `json:"param_id"`
	Value    int    `json:"value"`
}

func (s service) generalMessageProcessor(topicName, msgBody string) {
	var cont []flagMsg
	err := json.Unmarshal([]byte(msgBody), &cont)
	if err != nil {
		logrus.Debug("failed to deser msg body", "body", msgBody)
		return
	}
	m := make(map[string][]*model.Flag)
	for _, c := range cont {
		f := model.Flag{ChangeTime: time.Now()}
		if c.ParamID == 1 {
			if c.Value == 0 {
				f.Number = 1
				f.Value = true
			} else {
				f.Number = 2
				f.Value = true
			}
		} else if c.ParamID == 2 && c.Value > 11 {
			f.Number = 3
			f.Value = true
		}
		m[c.DeviceID] = append(m[c.DeviceID], &f)
	}
	for k, v := range m {
		s.r.ModifyFlags(context.Background(), k, v)
	}
}

func (s service) DeleteDevice(ctx context.Context, deviceName string) error {
	e := ServiceError{Method: "delete device"}
	err := s.r.DeleteDevice(ctx, deviceName)
	if err != nil {
		return e.Update("db query", err)
	}
	s.q.UnsubscribeDevices(ctx, deviceName)
	return nil
}

func (s service) GetDevice(ctx context.Context, deviceName string) (*model.Device, error) {
	e := ServiceError{Method: "get device"}
	dev, err := s.r.GetDevice(ctx, deviceName)
	if err == NoDeviceFound {
		return nil, nil
	}
	if err != nil {
		return &model.Device{}, e.Update("db query", err)
	}
	for _, f := range dev.Flags {
		f.Name = fmt.Sprintf("flag%d", f.Number)
	}
	s.acceptenceTest(&dev)
	return &dev, nil
}

func (s service) GetDevices(ctx context.Context) ([]model.Device, error) {
	e := ServiceError{Method: "get devices"}
	devs, err := s.r.GetDevices(ctx)
	if err != nil {
		return nil, e.Update("db query", err)
	}
	if devs == nil {
		devs = []model.Device{}
	}
	for i, d := range devs {
		s.acceptenceTest(&d)
		devs[i] = d
	}
	return devs, nil
}

func (s service) acceptenceTest(d *model.Device) {
	res := true
	for _, f := range d.Flags {
		if !f.Value || time.Since(f.ChangeTime).Hours() > 72 {
			res = false
			break
		}
	}
	d.AcceptenceResult = res
}
