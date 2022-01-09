package models

type RedReplay struct {
	Id        uint   `json:"id"`
	UserName  string `json:"username" gorm:"column:username"`
	Attribute string `json:"attribute" gorm:"column:attribute"`
	Operation string `json:"op" gorm:"column:op"`
	Value     string `json:"value" gorm:"column:value"`
}
