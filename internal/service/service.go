package service

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/cache"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"
	"github.com/rs/zerolog"
)

type ServiceHolder struct {
	Account *AccountService
	Item    *ItemService
}

func NewServiceHolder(rep *repository.Repository, itemCache *cache.ItemCache, log *zerolog.Logger) *ServiceHolder {
	return &ServiceHolder{
		Account: NewAccountService(rep, itemCache, log),
		Item:    NewItemService(itemCache, log),
	}
}
