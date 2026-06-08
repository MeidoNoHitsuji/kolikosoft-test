package service

import "github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"

type ItemService struct {
	rep *repository.Repository
}

func NewItemService(rep *repository.Repository) *ItemService {
	return &ItemService{
		rep: rep,
	}
}
