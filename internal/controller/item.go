package controller

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	itemSrv *service.ItemService
}

func NewItemController(backSrv *service.ItemService) *ItemController {
	return &ItemController{
		itemSrv: backSrv,
	}
}

func (c *ItemController) RegisterRoutes(router *gin.RouterGroup) {
	items := router.Group("/items")
	items.GET("/", c.getItems)
	items.POST("/buy", c.buyItem)
}

func (c *ItemController) getItems(ctx *gin.Context) {
	
}

func (c *ItemController) buyItem(ctx *gin.Context) {

}
