package controller

import (
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accSrv *service.AccountService
}

func NewAccountController(accSrv *service.AccountService) *AccountController {
	return &AccountController{
		accSrv: accSrv,
	}
}

func (c *AccountController) RegisterRoutes(router *gin.RouterGroup) {
	account := router.Group("/account")
	account.POST("/add-money", c.addMoney)
}

func (c *AccountController) addMoney(ctx *gin.Context) {

}
