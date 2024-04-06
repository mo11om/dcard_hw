package main

import (
	controllers "api/src/controllers"
	service "api/src/services"

	"github.com/gin-gonic/gin"
)

func main() {
	err := controllers.DBconnect()
	if err != nil {
		panic(err) // Handle error more gracefully in production
	}
	err = controllers.Init_redis()
	if err != nil {
		panic(err) // Handle error more gracefully in production
	}
	go func() {
		controllers.DBconnect()
		controllers.Init_redis()
	}()
	router := gin.Default()

	// Admin API - Create Ad

	router.POST("api/v1/ad", service.CreateAd)

	// // Public API - List Ads

	router.GET("api/v1/ad", service.ListAds)

	router.Run(":8080")
}
