package services

import (
	model "api/src/model"
	"fmt"
	"net/http"
	"strconv"

	controllers "api/src/controllers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Create Ad handler
func CreateAd(c *gin.Context) {
	var ad model.Ad
	if err := c.ShouldBindBodyWith(&ad, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate parameters
	if ad.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	if ad.StartAt.IsZero() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "StartAt is required"})
		return
	}
	if ad.EndAt.IsZero() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "EndAt is required"})
		return
	}
	if ad.StartAt.After(ad.EndAt) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "StartAt must be before EndAt"})
		return
	}

	// Save ad to database (replace with your actual logic)
	// ...
	controllers.Create_ad(ad)

	// Save conditions to database (replace with your actual logic)
	// ...
	// controllers.Create_condition(ad.Conditions[0], 3)

	c.Status(http.StatusCreated)
}

// List Ads handler
func ListAds(c *gin.Context) {
	//check
	var serach_condition model.Search_Condition
	age, err := strconv.Atoi(c.Query("age"))
	if err != nil {
		age = 0
	}
	gender := c.Query("gender")
	country := c.Query("country")
	platform := c.Query("platform")

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 1 || offset > 100 {
		offset = 5
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "offset must be num or not in 1-100"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 5
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit must be num or not in 1-100"})

	}
	//fetch params into data structure
	serach_condition.Age = age
	serach_condition.Country = country
	serach_condition.Gender = gender
	serach_condition.Platform = platform
	serach_condition.Limit = limit
	serach_condition.Offset = offset

	// Fetch ads from database based on filters, pagination, and active status (replace with your actual logic)

	ads, err := controllers.Find_ad(serach_condition)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(500)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"items": ads,
		})
	}

	// Replace "ads" with actual fetched data
}
