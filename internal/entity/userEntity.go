package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey"`
	Username  string    `gorm:"column:username;uniqueIndex; not null" type:"varchar(255)"`
	Name      string    `gorm:"column:name; not null" type:"varchar(255)"`
	Email     string    `gorm:"column:email;uniqueIndex; not null" type:"varchar(255)"`
	Password  string    `gorm:"column:password; not null" type:"varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
