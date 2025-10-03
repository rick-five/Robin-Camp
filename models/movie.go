package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Movie struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title       string    `json:"title" gorm:"uniqueIndex;not null"`
	ReleaseDate time.Time `json:"releaseDate" gorm:"not null"`
	Genre       string    `json:"genre" gorm:"not null"`
	Distributor string    `json:"distributor,omitempty"`
	Budget      int64     `json:"budget,omitempty"`
	MPARating   string    `json:"mpaRating,omitempty"`
	BoxOffice   *BoxOffice `json:"boxOffice" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateTime"`
}

type BoxOffice struct {
	Revenue      Revenue `json:"revenue"`
	Currency     string  `json:"currency"`
	Source       string  `json:"source"`
	LastUpdated  string  `json:"lastUpdated"`
}

type Revenue struct {
	Worldwide        int64 `json:"worldwide"`
	OpeningWeekendUsa int64 `json:"openingWeekendUsa,omitempty"`
}

// BeforeCreate 会在创建记录前调用
func (m *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果ID为空，则生成一个新的UUID
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}