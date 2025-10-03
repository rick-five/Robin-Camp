package models

import (
	"time"
)

type Rating struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MovieTitle string    `json:"movieTitle" gorm:"index;not null"`
	RaterID    string    `json:"-" gorm:"index;not null"` // 不在JSON中暴露
	Rating     float64   `json:"rating" gorm:"type:decimal(2,1);check:rating >= 0.5 AND rating <= 5.0 AND rating * 2 = ROUND(rating * 2);not null"`
	CreatedAt  time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"-" gorm:"autoUpdateTime"`
}

type RatingSubmit struct {
	Rating float64 `json:"rating" binding:"required"`
}

type RatingResult struct {
	MovieTitle string  `json:"movieTitle"`
	RaterID    string  `json:"raterId"`
	Rating     float64 `json:"rating"`
}

type RatingAggregate struct {
	Average float64 `json:"average"`
	Count   int64   `json:"count"`
}
