package request

import (
	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/gin-gonic/gin"
)

type AccountAddMoney struct {
	ID    int64 `json:"-"`
	Value int64 `json:"value" binding:"required,min=1"`
}

func (i *AccountAddMoney) Validate(ginCtx *gin.Context) (*model.AccountAddMoney, error) {
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

func (i *AccountAddMoney) ToEntity() *model.AccountAddMoney {
	return &model.AccountAddMoney{
		ID:    i.ID,
		Value: i.Value,
	}
}
