package ds

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"type:varchar(40)"`
	Password string `gorm:"type:varchar(40)"`
}
