package service

import "github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"

type ServiceHolder struct {
	Account *AccountService
	Item    *ItemService
}

func NewServiceHolder(rep *repository.Repository) *ServiceHolder {
	return &ServiceHolder{
		Account: NewAccountService(rep),
		Item:    NewItemService(rep),
	}
}
