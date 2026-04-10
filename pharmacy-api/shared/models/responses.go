package models

type GoodsResponse struct {
	ID                   uint     `json:"ID"`
	Name                 string   `json:"Name"`
	Image                string   `json:"Image"`
	Producer             string   `json:"Producer"`
	IsInStock            bool     `json:"IsInStock"`
	Tags                 []string `json:"Tags"`
	Instruction          string   `json:"Instruction"`
	Description          string   `json:"Description"`
	IsPrescriptionNeeded bool     `json:"IsPrescriptionNeeded"`
	Price                uint     `json:"Price"`
}

type ScheduleResponse struct {
	ID        uint   `json:"ID"`
	Date      string `json:"Date"`
	TimeStart string `json:"TimeStart"`
	TimeEnd   string `json:"TimeEnd"`
	IsOpened  bool   `json:"IsOpened"`
}

type AnnouncementResponse struct {
	ID       uint   `json:"ID"`
	DateTime string `json:"DateTime"`
	From     string `json:"From"`
	Announce string `json:"Announce"`
}

type OrderResponse struct {
	ID    uint          `json:"ID"`
	Name  string        `json:"Name"`
	Email string        `json:"Email,omitempty"`
	Phone string        `json:"Phone"`
	Items []OrderedItem `json:"Items"`
}

type OrderedItem struct {
	ID          uint   `json:"ID"`
	Name        string `json:"Name"`
	Image       string `json:"Image"`
	Description string `json:"Description"`
	Price       uint   `json:"Price"`
	Quantity    uint   `json:"Quantity"`
}

type GoodsUpdateRequest struct {
	ID                   uint   `json:"ID"`
	Name                 string `json:"Name"`
	Instruction          string `json:"Instruction"`
	Description          string `json:"Description"`
	IsPrescriptionNeeded bool   `json:"IsPrescriptionNeeded"`
	IsInStock            bool   `json:"IsInStock"`
	Price                uint   `json:"Price"`
}

type PromoItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Image       string `json:"image"`
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
