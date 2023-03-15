package entity

import "time"

type Server struct {
	Id        int        `json:"id" gorm:"column:id;type:uuid;uuid_generate_v4();primary_key"`
	Name      string     `json:"name" gorm:"column:name;type:varchar(255);not null;uniqueIndex"`
	Ipv4      string     `json:"ipv4" gorm:"column:ipv4;type:varchar(15);not null;uniqueIndex"`
	User      string     `json:"user" gorm:"column:user;type:varchar(50)"`
	Password  string     `json:"password" gorm:"column:password;type:varchar(100)"`
	Status    string     `json:"status" gorm:"column:status;type:varchar(50)"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (Server) TableName() string { return "servers" }
