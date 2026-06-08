package controller

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	srv *service.ItemService
}

func NewItemController(backSrv *service.ItemService) *ItemController {
	return &ItemController{
		srv: backSrv,
	}
}

func (c *ItemController) RegisterRoutes(router *gin.RouterGroup) {
	items := router.Group("/items")
	items.GET("/", c.getItems)
}

func (c *ItemController) getItems(ctx *gin.Context) {

}
