package router

import (
	"ginchat/service"
	"net/http"

	docs "ginchat/docs"

	"ginchat/common"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func handleCors(router *gin.Engine) {
	// CORS 中间件，可以根据需要自行配置
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

}

func Router() *gin.Engine {
	r := gin.Default()

	handleCors(r)

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/ping", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/register", service.RegisterUser)
	r.POST("/user/login", service.GetUserByName)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/updatePassword", service.UpdatePassword)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)

	r.POST("/user/getFriends", common.TokenExtractionMiddleware(), service.SearchFriends)
	return r
}
