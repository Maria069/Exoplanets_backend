package ds

type Likes struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;uniqueIndex:idx_render_server_unit_user"`
	ExoplanetID uint      `gorm:"not null;uniqueIndex:idx_render_server_unit_user"`
	User        User      `gorm:"foreignKey:UserID"`
	Exoplanet   Exoplanet `gorm:"foreignKey:ExoplanetID"`
}
