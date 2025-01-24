package entity

import "time"

type User struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Username  string    `gorm:"column:username;unique" type:"varchar(255)"`
	Name      string    `gorm:"column:name" type:"varchar(255)"`
	Email     string    `gorm:"column:email;unique" type:"varchar(255)"`
	Password  string    `gorm:"column:password" type:"varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}
