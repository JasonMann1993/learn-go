package httptest_demo

import (
	"fmt"
	"net/Http"

	"github.com/gin-gonic/gin"
)

type Param struct {
	Name string `json:"name"`
}

func helloHandler(c *gin.Context){
	var p Param
}