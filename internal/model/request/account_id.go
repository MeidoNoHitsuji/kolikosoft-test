package request

import (
	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/gin-gonic/gin"
)

type AccountId struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (i *AccountId) Validate(ginCtx *gin.Context) (*model.AccountId, error) {
	err := ginCtx.ShouldBindUri(&i)
	if err != nil {
		return nil, backErr.ErrValidate
	}

	return i.ToEntity(), nil
}

func (i *AccountId) ToEntity() *model.AccountId {
	return &model.AccountId{
		ID: i.ID,
	}
}
