package models

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(100);unique;not null"`
	Password    string `gorm:"type:varchar(255);not null"`
	IsModerator bool   `gorm:"type:boolean;default:false"`
}

func (User) TableName() string { return "users" }
