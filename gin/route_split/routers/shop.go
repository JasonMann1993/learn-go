package routers

import "github.com/gin-gonic/gin"

func goodsHandler(c *gin.Context) {

}

func checkoutHandler(c *gin.Context) {

}

func LoadShop(e *gin.Engine)  {
	e.GET("/goods", goodsHandler)
	e.GET("/checkout", checkoutHandler)
}