package model

import (
	"time"
)

// Ad data structure

// Redis connection (replace with your actual connection logic)
// var redisClient *redis.Client
type Ad struct {
	// ID         int         `json:"id"`
	Title      string      `json:"title"`
	StartAt    time.Time   `json:"startAt"`
	EndAt      time.Time   `json:"endAt"`
	Conditions []Condition `json:"conditions"`
}

// Condition data structure
type Condition struct {
	AgeStart int      `json:"ageStart"`
	AgeEnd   int      `json:"ageEnd"`
	Gender   []string `json:"gender"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
}

// Condition to serach data  structure
type Search_Condition struct {
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Country  string `json:"country"`
	Platform string `json:"platform"`
	Limit    int
	Offset   int
}

type Result struct {
	Title string    `json:"title" gorm:"title"`
	EndAt time.Time `json:"endAt" gorm:"end_at"`
}

var All_platform = []string{"android", "ios", "web"}
var All_gender = []string{"M", "F"}
