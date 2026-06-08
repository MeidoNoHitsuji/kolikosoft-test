package controller

import (
	"errors"
	"net/http"

	errors2 "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model/response"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/gin-gonic/gin"
)

type ControllerInterface interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type ControllerHolder struct {
	Controllers []ControllerInterface
}

func NewControllerHolder(srvHodler *service.ServiceHolder) *ControllerHolder {
	return &ControllerHolder{
		Controllers: []ControllerInterface{
			NewAccountController(srvHodler.Account),
			NewItemController(srvHodler.Item),
		},
	}
}

func errorHandler(ginCtx *gin.Context, err error) {
	var errStr *errors2.ErrorStruct

	if errors.As(err, &errStr) {
		ginCtx.JSON(errStr.Code, errStr.ToResponse())
		return
	}

	ginCtx.JSON(http.StatusInternalServerError, response.NewError(err))
}
