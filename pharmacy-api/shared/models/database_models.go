package models

import (
	"database/sql"
	"time"
)

// MYSQL or GORM moment idk, on Postgres all alright
type WorkTime struct {
	ID        uint           `gorm:"primaryKey"`
	Date      sql.NullString `gorm:"type:DATE; not null "`
	IsOpened  bool           `gorm:"not null"`
	TimeStart sql.NullString `gorm:"type:TIME"`
	TimeEnd   sql.NullString `gorm:"type:TIME"`
}

type Producer struct {
	ID           uint    `gorm:"primaryKey"`
	ProducerName string  `gorm:"size:40; not null"`
	Goods        []Goods `gorm:"foreignKey:ProducerID"`
}

type Tag struct {
	ID      uint    `gorm:"primaryKey"`
	TagName string  `gorm:"size:20; not null"`
	Goods   []Goods `gorm:"many2many:goods_tags"`
}

type GoodsOrders struct {
	OrderID  uint  `gorm:"primaryKey"`
	GoodsID  uint  `gorm:"primaryKey"`
	Quantity uint  `gorm:"not null"`
	Order    Order `gorm:"foreignKey:OrderID"`
	Goods    Goods `gorm:"foreignKey:GoodsID"`
}

type Order struct {
	ID          uint          `gorm:"primaryKey"`
	ClientFIO   string        `gorm:"size:30; not null"`
	ClientEmail string        `gorm:"size:30"`
	ClientPhone string        `gorm:"size:18; not null"`
	Goods       []GoodsOrders `gorm:"foreignKey:OrderID; constraint:OnDelete:CASCADE"`
}

type Goods struct {
	ID                   uint          `gorm:"primaryKey"`
	Name                 string        `gorm:"size:64; not null"`
	Image                string        `gorm:"size:64; not null"`
	ProducerID           uint          `gorm:"index; not null"`
	Producer             Producer      `gorm:"foreignKey:ProducerID; not null"`
	IsInStock            bool          `gorm:"not null"`
	Tags                 []Tag         `gorm:"many2many:goods_tags; not null"`
	Orders               []GoodsOrders `gorm:"foreignKey:GoodsID"`
	Instruction          string        `gorm:"size:1024"`
	Description          string        `gorm:"size:1024"`
	IsPrescriptionNeeded bool          `gorm:"not null"`
	Price                uint          `gorm:"not null"`
}

type Announcement struct {
	ID       uint      `gorm:"primaryKey"`
	DateTime time.Time `gorm:"not null"`
	From     string    `gorm:"size:64; not null"`
	Announce string    `gorm:"size:2048; not null"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Login        string `gorm:"size:40; not null; unique; index"`
	UserName     string `gorm:"size:40; not null"`
	RoleID       uint   `gorm:"index; not null"`
	Role         Role   `gorm:"foreignKey:RoleID; not null"`
	PasswordHash []byte `gorm:"size:128; not null"`
}

type Role struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"size:20; not null"`
	User        []User       `gorm:"foreignKey:RoleID"`
	Permissions []Permission `gorm:"many2many:role_permission; not null; constraint:OnDelete:CASCADE"`
}

type Permission struct {
	ID     uint   `gorm:"primaryKey"`
	Action string `gorm:"size:20; not null"`
	Role   []Role `gorm:"many2many:role_permission"`
}
