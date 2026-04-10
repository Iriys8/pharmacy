package models

import (
	"github.com/golang-jwt/jwt/v4"
)

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

type RefreshClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
