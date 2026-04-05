package models

type GoodsUpdateRequest struct {
	ID                   uint   `json:"ID"`
	Name                 string `json:"Name"`
	Instruction          string `json:"Instruction"`
	Description          string `json:"Description"`
	IsPrescriptionNeeded bool   `json:"IsPrescriptionNeeded"`
	IsInStock            bool   `json:"IsInStock"`
	Price                uint   `json:"Price"`
}

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
