package models

type UserInfo struct {
	Id         uint   `json:"id"`
	UserName   string `json:"username" gorm:"column:username"`
	Mail       string `json:"mail" gorm:"column:mail"`
	Name       string `json:"name" gorm:"column:name"`
	Department string `json:"department" gorm:"column:department"`
	WorkPhone  string `json:"workphone" gorm:"column:workphone"`
	HomePhone  string `json:"homephone" gorm:"column:homephone"`
	Mobile     string `json:"mobile" gorm:"column:mobile"`
}
