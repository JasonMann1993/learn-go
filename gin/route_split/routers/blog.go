package routers

import "github.com/gin-gonic/gin"

func postHandler(c *gin.Context) {

}

func commentHandler(c *gin.Context) {

}

func LoadBlog(e *gin.Engine)  {
	e.GET("/post", postHandler)
	e.GET("/comment", commentHandler)
}