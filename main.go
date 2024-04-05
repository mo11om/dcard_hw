package main

import (
	database "api/src/controllers"
	service "api/src/services"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Admin API - Create Ad

	router.POST("api/v1/ad", service.CreateAd)

	// // Public API - List Ads
	go func() {
		database.DBconnect()
	}()
	router.GET("api/v1/ad", service.ListAds)

	router.Run(":8080")
}
