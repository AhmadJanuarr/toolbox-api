package routes

import (
	"toolkits/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	//configuration router dengan gin default
	router := gin.Default()

	// group endpoint
	api := router.Group("/api/v1")

	// endpoint image
	api.POST("/image/convert", handlers.ConvertImage)
	api.POST("/image/compress-image", handlers.CompressionImage)
	api.POST("/image/resize-image", handlers.ResizeImage)

	//enpoint downloader
	return router
}
