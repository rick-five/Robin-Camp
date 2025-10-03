package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"robin-camp/models"
	"robin-camp/utils"
)

type MovieHandler struct {
	DB *gorm.DB
}

type MovieCreateRequest struct {
	Title       string    `json:"title" binding:"required"`
	ReleaseDate time.Time `json:"releaseDate" binding:"required" time_format:"2006-01-02"`
	Genre       string    `json:"genre" binding:"required"`
	Distributor string    `json:"distributor,omitempty"`
	Budget      int64     `json:"budget,omitempty"`
	MPARating   string    `json:"mpaRating,omitempty"`
}

func NewMovieHandler(db *gorm.DB) *MovieHandler {
	return &MovieHandler{DB: db}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req MovieCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建电影对象
	movie := models.Movie{
		Title:       req.Title,
		ReleaseDate: req.ReleaseDate,
		Genre:       req.Genre,
		Distributor: req.Distributor,
		Budget:      req.Budget,
		MPARating:   req.MPARating,
	}

	// 保存到数据库
	if err := h.DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	// 调用票房API获取数据
	boxOfficeData, err := utils.FetchBoxOfficeData(req.Title)
	if err == nil {
		// 如果成功获取票房数据，更新电影记录
		boxOffice := &models.BoxOffice{
			Revenue: models.Revenue{
				Worldwide:         boxOfficeData.Revenue.Worldwide,
				OpeningWeekendUsa: boxOfficeData.Revenue.OpeningWeekendUsa,
			},
			Currency:    "USD",
			Source:      "BoxOfficeAPI",
			LastUpdated: time.Now().Format(time.RFC3339),
		}

		// 用户提供的值优先于票房API数据
		if movie.Distributor == "" {
			movie.Distributor = boxOfficeData.Distributor
		}
		if movie.Budget == 0 {
			movie.Budget = boxOfficeData.Budget
		}
		if movie.MPARating == "" {
			movie.MPARating = boxOfficeData.MPARating
		}
		movie.BoxOffice = boxOffice

		// 更新数据库中的电影记录
		h.DB.Save(&movie)
	} else {
		// 如果票房API调用失败，将boxOffice设置为null
		movie.BoxOffice = nil
		h.DB.Save(&movie)
	}

	// 返回创建的电影
	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	var movies []models.Movie

	query := h.DB.Model(&models.Movie{})

	// 处理搜索参数
	if q := c.Query("q"); q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}

	if year := c.Query("year"); year != "" {
		query = query.Where("EXTRACT(YEAR FROM release_date) = ?", year)
	}

	if genre := c.Query("genre"); genre != "" {
		query = query.Where("genre = ?", genre)
	}

	if distributor := c.Query("distributor"); distributor != "" {
		query = query.Where("distributor = ?", distributor)
	}

	if budget := c.Query("budget"); budget != "" {
		query = query.Where("budget <= ?", budget)
	}

	if mpaRating := c.Query("mpaRating"); mpaRating != "" {
		query = query.Where("mpa_rating = ?", mpaRating)
	}

	// 处理分页
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		// 这里简化处理，实际应该解析为整数并验证
	}

	query.Limit(limit).Find(&movies)

	// 构造响应
	response := gin.H{
		"items": movies,
	}

	c.JSON(http.StatusOK, response)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie

	if err := h.DB.First(&movie, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}
