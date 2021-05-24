package endpoints

import (
	"context"
	"github.com/QuanAnh/rabbitmq-go/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetUserById				endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		GetUserById: makeGetUserById(s),
	}
}

func makeGetUserById(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		u := request.(UserRequest)
		return s.GetUserById(ctx, u.Id)
	}
}
