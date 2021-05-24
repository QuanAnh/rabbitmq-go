package service

import (
	"context"
	"github.com/jinzhu/gorm"
)

type repo struct {
	db		*gorm.DB
}

func NewRepo(connection *gorm.DB) *repo {
	return &repo{
		db: connection,
	}
}

type Repo interface {
	GetUserById(ctx context.Context, id string) (interface{}, error)
}

func (r *repo) GetUserById(ctx context.Context, id string) (interface{}, error) {
	return "ok", nil
}