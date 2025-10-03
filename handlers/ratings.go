package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"robin-camp/models"
)

type RatingHandler struct {
	DB *gorm.DB
}

func NewRatingHandler(db *gorm.DB) *RatingHandler {
	return &RatingHandler{DB: db}
}

func (h *RatingHandler) SubmitRating(c *gin.Context) {
	title := c.Param("title")
	raterID := c.GetHeader("X-Rater-Id")

	// 检查raterID是否存在
	if raterID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing X-Rater-Id header"})
		return
	}

	// 检查电影是否存在
	var movie models.Movie
	if err := h.DB.Where("title = ?", title).First(&movie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// 解析请求体
	var req models.RatingSubmit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查评分值是否有效
	validRatings := map[float64]bool{
		0.5: true, 1.0: true, 1.5: true, 2.0: true, 2.5: true,
		3.0: true, 3.5: true, 4.0: true, 4.5: true, 5.0: true,
	}

	if !validRatings[req.Rating] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid rating value"})
		return
	}

	// 创建或更新评分
	rating := models.Rating{
		MovieTitle: title,
		RaterID:    raterID,
		Rating:     req.Rating,
	}

	// 检查评分是否已存在
	var existingRating models.Rating
	result := h.DB.Where("movie_title = ? AND rater_id = ?", title, raterID).First(&existingRating)

	if result.Error != nil {
		// 评分不存在，创建新评分
		if err := h.DB.Create(&rating).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
			return
		}
		c.JSON(http.StatusCreated, models.RatingResult{
			MovieTitle: rating.MovieTitle,
			RaterID:    rating.RaterID,
			Rating:     rating.Rating,
		})
	} else {
		// 评分已存在，更新评分
		existingRating.Rating = req.Rating
		if err := h.DB.Save(&existingRating).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rating"})
			return
		}
		c.JSON(http.StatusOK, models.RatingResult{
			MovieTitle: existingRating.MovieTitle,
			RaterID:    existingRating.RaterID,
			Rating:     existingRating.Rating,
		})
	}
}

func (h *RatingHandler) GetRatingAggregate(c *gin.Context) {
	title := c.Param("title")

	// 检查电影是否存在
	var movie models.Movie
	if err := h.DB.Where("title = ?", title).First(&movie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// 计算评分聚合
	var result struct {
		Average float64
		Count   int64
	}

	row := h.DB.Model(&models.Rating{}).
		Select("COALESCE(AVG(rating), 0) as average, COUNT(*) as count").
		Where("movie_title = ?", title).
		Row()

	if err := row.Scan(&result.Average, &result.Count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate rating aggregate"})
		return
	}

	// 四舍五入到1位小数
	average := float64(int(result.Average*10+0.5)) / 10

	c.JSON(http.StatusOK, models.RatingAggregate{
		Average: average,
		Count:   result.Count,
	})
}
