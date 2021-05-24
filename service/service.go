package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type service struct {
	repository				*repo
	logger					log.Logger
}

func NewBasicService(repo *repo, logger log.Logger) *service {
	return &service{
		repository: repo,
		logger: logger,
	}
}

type Service interface {
	GetUserById(ctx context.Context, id string) (interface{}, error)
}

func (s *service) GetUserById(ctx context.Context, id string) (interface{}, error) {
	user, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		return "", err
	}
	return user, nil
}
