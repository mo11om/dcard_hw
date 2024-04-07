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
	// Bind incoming request body to an Ad struct
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

	// Save ad to database
	// Save conditions to database
	// ...
	if ad.Conditions[0].Country == nil {
		ad.Conditions[0].Country = []string{"ALL"}
	}
	if ad.Conditions[0].Gender == nil {
		ad.Conditions[0].Gender = model.All_gender
	}
	if ad.Conditions[0].Platform == nil {
		ad.Conditions[0].Platform = model.All_platform
	}
	fmt.Println(ad)

	controllers.Create_ad(ad)

	c.Status(http.StatusCreated)
}

// / ListAds handles listing ads based on search filters and pagination.
func ListAds(c *gin.Context) {
	// Extract query parameters for filtering and pagination
	var serach_condition model.Search_Condition
	age, err := strconv.Atoi(c.Query("age"))
	if err != nil {
		age = 0
	}
	gender := c.Query("gender")
	country := c.Query("country")
	platform := c.Query("platform")
	// Validate offset and limit parameters (ensure values are within range)
	offset, err := strconv.Atoi(c.Query("offset"))

	if err != nil || offset < 1 || offset > 100 {
		offset = 5
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 5
	}

	// Populate search condition struct with extracted parameters
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
