package controller

import (
	"net/http"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	srv *service.ItemService
}

func NewItemController(srv *service.ItemService) *ItemController {
	return &ItemController{
		srv: srv,
	}
}

func (c *ItemController) RegisterRoutes(router *gin.RouterGroup) {
	items := router.Group("/items")
	items.GET("/", c.getItems)
}

func (c *ItemController) getItems(ginCtx *gin.Context) {
	ctx := ginCtx.Request.Context()
	items, err := c.srv.GetItems(ctx)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, items)
}
