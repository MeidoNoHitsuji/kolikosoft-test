package service

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"
	"github.com/rs/zerolog"
)

type ItemService struct {
	rep *repository.Repository
	log *zerolog.Logger
}

func NewItemService(rep *repository.Repository, log *zerolog.Logger) *ItemService {
	return &ItemService{
		rep: rep,
		log: log,
	}
}
