package main

import (
	controllers "api/src/controllers"
	service "api/src/services"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Admin API - Create Ad

	router.POST("api/v1/ad", service.CreateAd)

	// // Public API - List Ads
	go func() {
		controllers.DBconnect()
		controllers.Init_redis()
	}()
	router.GET("api/v1/ad", service.ListAds)

	router.Run(":8080")
}
