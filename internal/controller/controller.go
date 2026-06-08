package controller

import (
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
