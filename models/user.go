package models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName" gorm:"column:userName"`
	Email    string `json:"email" gorm:"unique;column:email"`
	Password string `json:"password" gorm:"column:password"`
}
