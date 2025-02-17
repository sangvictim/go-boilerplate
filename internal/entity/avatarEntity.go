package entity

import (
	"gorm.io/gorm"
)

type Avatar struct {
	gorm.Model
	UserID       string
	OriginalName string `gorm:"column:original_name; not null" type:"varchar(255)"`
	Key          string `gorm:"column:key; not null" type:"varchar(255)"`
	Size         int64  `gorm:"column:size; not null"`
	MimeType     string `gorm:"column:mime_type; not null" type:"varchar(255)"`
	Visibility   string `gorm:"column:visibility; not null" enum:"PUBLIC,PRIVATE,ShARED" default:"PRIVATE" type:"varchar(255)"`
}

func (u *Avatar) TableName() string {
	return "avatar"
}
