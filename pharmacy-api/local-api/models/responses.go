package models

type UserUpdateRequest struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"size:40; not null; unique; index"`
	UserName string `gorm:"size:40; not null"`
	RoleID   uint   `gorm:"index; not null"`
	Role     Role   `gorm:"foreignKey:RoleID; not null"`
	Password string `gorm:"size:128; not null"`
}

type LogsResponse struct {
	Name string `json:"Name"`
}
