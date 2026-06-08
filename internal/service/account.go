package service

import (
	"context"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/cache"
	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"
	"github.com/rs/zerolog"
)

type AccountService struct {
	rep       *repository.Repository
	itemCache *cache.ItemCache
	log       *zerolog.Logger
}

func NewAccountService(rep *repository.Repository, itemCache *cache.ItemCache, log *zerolog.Logger) *AccountService {
	return &AccountService{
		rep:       rep,
		itemCache: itemCache,
		log:       log,
	}
}

func (s *AccountService) GetAccountByID(ctx context.Context, id int64) (*model.Account, error) {
	acc, err := s.rep.GetAccountByID(ctx, id)
	if err != nil {
		s.log.Err(err).Msg("Не удалось найти аккаунт в базе данных")
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccountHistoriesByAccountID(ctx context.Context, id int64) ([]model.AccountHistory, error) {
	history, err := s.rep.GetAccountHistoriesByAccountID(ctx, id)
	if err != nil {
		s.log.Err(err).Msg("Не удалось найти аккаунт в базе данных")
		return nil, err
	}

	return history, nil
}

func (s *AccountService) AddMoney(ctx context.Context, data *model.AccountAddMoney) (*model.Account, error) {
	return s.changeBalance(ctx, data, "Пополнение баланса")
}

func (s *AccountService) BuyItem(ctx context.Context, data *model.AccountBuyItem) (*model.Account, error) {
	item, err := s.itemCache.GetItem(ctx, data.HashMarketName)
	if err != nil {
		s.log.Err(err).Str("hash_market_name", data.HashMarketName).Msg("Не удалось найти предмет в кеше")
		return nil, err
	}

	if item.SuggestedPrice == nil {
		return nil, backErr.ErrItemPriceEmpty
	}

	value := int64(*item.SuggestedPrice * 100)
	txData := &model.AccountAddMoney{
		ID:    data.ID,
		Value: -value,
	}

	return s.changeBalance(ctx, txData, "Покупка предмета: "+data.HashMarketName)
}

func (s *AccountService) changeBalance(ctx context.Context, data *model.AccountAddMoney, comment string) (*model.Account, error) {
	tx, err := s.rep.BeginTx(ctx)
	if err != nil {
		s.log.Err(err).Msg("Не удалось начать транзакцию")
		return nil, backErr.ErrInternalServer
	}

	result, err := s.rep.TxUpdateBalance(ctx, tx, data)
	if err != nil {
		if errTx := s.rep.Rollback(tx); errTx != nil {
			s.log.Err(errTx).Msg("Не удалось откатить транзакцию")
			return nil, backErr.ErrInternalServer
		}
		return nil, err
	}

	result.Comment = comment

	err = s.rep.TxAddHistory(ctx, tx, result)
	if err != nil {
		if errTx := s.rep.Rollback(tx); errTx != nil {
			s.log.Err(errTx).Msg("Не удалось откатить транзакцию")
			return nil, backErr.ErrInternalServer
		}
		return nil, err
	}

	err = s.rep.Commit(tx)
	if err != nil {
		s.log.Err(err).Msg("Не удалось накатить транзакцию")
		if errTx := s.rep.Rollback(tx); errTx != nil {
			s.log.Err(errTx).Msg("Не удалось откатить транзакцию")
			return nil, backErr.ErrInternalServer
		}
		return nil, backErr.ErrInternalServer
	}

	return &model.Account{
		ID:      result.AccountID,
		Balance: result.NewBalance,
	}, nil
}
