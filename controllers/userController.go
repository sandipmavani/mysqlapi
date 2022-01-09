package controllers

import (
	"fmt"
	"mysqlapi/database"
	"mysqlapi/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type CreateResponse struct {
	UserName string `json:"username"`
	Id       uint   `json:"id"`
	Mail     string `json:"mail`
}

func CreateUser(c *fiber.Ctx) error {

	var data map[string]string

	if error := c.BodyParser(&data); error != nil {
		return error
	}

	var checkuser models.UserInfo

	database.DB.Table("userinfo").Where("username = ?", data["username"]).First(&checkuser)

	if checkuser.UserName == data["username"] {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "username already exist",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	fmt.Println(password)
	userinfo := models.UserInfo{
		Name:     data["name"],
		Mail:     data["mail"],
		UserName: data["username"],
		Mobile:   data["mobile"],
	}

	database.DB.Table("userinfo").Create(&userinfo)

	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
	// 	c.Status(fiber.StatusBadRequest)
	// 	return c.JSON(fiber.Map{
	// 		"message": "incorrect password",
	// 	})
	// }

	var tmp CreateResponse
	tmp.UserName = userinfo.UserName
	tmp.Id = userinfo.Id
	tmp.Mail = userinfo.Mail
	return c.JSON(fiber.Map{
		"message": "success",
		"user":    tmp,
	})

}

func GetAllUser(c *fiber.Ctx) error {

	var userInfoList []models.UserInfo
	database.DB.Table("userinfo").Find(&userInfoList)

	return c.JSON(userInfoList)

}
