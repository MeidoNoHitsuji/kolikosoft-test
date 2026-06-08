package controller

import (
	"net/http"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model/request"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	srv *service.AccountService
}

func NewAccountController(accSrv *service.AccountService) *AccountController {
	return &AccountController{
		srv: accSrv,
	}
}

func (c *AccountController) RegisterRoutes(router *gin.RouterGroup) {
	accountGroup := router.Group("/account/:id")
	accountGroup.GET("/", c.getInfo)
	accountGroup.GET("/history", c.getHistory)
	accountGroup.POST("/add-money", c.addMoney)
	accountGroup.POST("/buy-item", c.buyItem)
}

func (c *AccountController) getInfo(ginCtx *gin.Context) {
	var r request.AccountId

	data, err := r.Validate(ginCtx)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ctx := ginCtx.Request.Context()
	account, err := c.srv.GetAccountByID(ctx, data.ID)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, account)
}

func (c *AccountController) getHistory(ginCtx *gin.Context) {
	var r request.AccountId

	data, err := r.Validate(ginCtx)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ctx := ginCtx.Request.Context()
	history, err := c.srv.GetAccountHistoriesByAccountID(ctx, data.ID)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, history)
}

func (c *AccountController) addMoney(ginCtx *gin.Context) {
	var r request.AccountAddMoney

	data, err := r.Validate(ginCtx)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ctx := ginCtx.Request.Context()
	account, err := c.srv.AddMoney(ctx, data)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, account)
}

func (c *AccountController) buyItem(ginCtx *gin.Context) {
	var r request.AccountBuyItem

	data, err := r.Validate(ginCtx)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ctx := ginCtx.Request.Context()
	account, err := c.srv.BuyItem(ctx, data)
	if err != nil {
		errorHandler(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, account)
}
