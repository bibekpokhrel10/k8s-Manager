package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wordpress.com/internal/controllers"
)

func RouterHandle() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"API Health": "OK"})
	})
	r.POST("/create", controllers.AppCreate)
	r.DELETE("/delete", controllers.AppDelete)
	r.GET("/list", controllers.ListAll)
	r.GET("/detail", controllers.GetDetails)
	r.PUT("/update", controllers.AppUpdate)
	r.Run(":4000")
}
