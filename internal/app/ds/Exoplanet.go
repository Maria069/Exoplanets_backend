package ds

import "time"

type Exoplanet struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(30)" json:"title"`
	Status      string    `gorm:"type:varchar(30)" json:"status"`
	Distance    int       `json:"distance"`
	Weight      float64   `gorm:"type:numeric(6,2)" json:"weight"`
	Description string    `gorm:"type:varchar(300)" json:"description"`
	Image       string    `gorm:"type:varchar(100)" json:"image"`
	Video       string    `gorm:"type:varchar(100)" json:"video"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	CreatorID   uint      `gorm:"index" json:"-"`
	Creator     User      `gorm:"foreignKey:CreatorID" json:"-"`
	FormedAt    time.Time `gorm:"autoUpdateTime" json:"-"`
}
