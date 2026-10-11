package ds

import "time"

type Exoplanet struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(30)"`
	Status      string `gorm:"type:varchar(30)"`
	Distance    int
	Weight      float64   `gorm:"type:numeric(6,2)"`
	Description string    `gorm:"type:varchar(300)"`
	ImagePath   string    `gorm:"type:varchar(100)"`
	VideoPath   string    `gorm:"type:varchar(100)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	CreatorID   uint      `gorm:"index"`
	Creator     User      `gorm:"foreignKey:CreatorID"`
	FormedAt    time.Time `gorm:"autoUpdateTime"`
}
