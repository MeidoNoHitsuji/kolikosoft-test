package service

import "github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"

type AccountService struct {
	rep *repository.Repository
}

func NewAccountService(rep *repository.Repository) *AccountService {
	return &AccountService{
		rep: rep,
	}
}
