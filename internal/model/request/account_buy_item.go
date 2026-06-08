package request

import (
	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/gin-gonic/gin"
)

type AccountBuyItem struct {
	ID             int64  `json:"-"`
	HashMarketName string `json:"hash_market_name" binding:"required"`
}

func (i *AccountBuyItem) Validate(ginCtx *gin.Context) (*model.AccountBuyItem, error) {
	var accId AccountId
	acc, err := accId.Validate(ginCtx)
	if err != nil {
		return nil, backErr.ErrValidate
	}

	err = ginCtx.ShouldBindJSON(&i)
	if err != nil {
		return nil, backErr.ErrValidate
	}

	i.ID = acc.ID
	return i.ToEntity(), nil
}

func (i *AccountBuyItem) ToEntity() *model.AccountBuyItem {
	return &model.AccountBuyItem{
		ID:             i.ID,
		HashMarketName: i.HashMarketName,
	}
}
