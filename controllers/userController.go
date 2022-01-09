package controllers

import (
	"fmt"
	"log"
	"mysqlapi/database"
	"mysqlapi/models"

	"github.com/gofiber/fiber/v2"
	"github.com/simia-tech/crypt"
	"github.com/syyongx/php2go"
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
	if value, ok := data["password"]; ok {
		updatePassword(userinfo.UserName, value)
	}

	if value, ok := data["sessionAllowed"]; ok {
		upsertOperation("radcheck", "Max-All-Session", value, userinfo.UserName)
	}

	if value, ok := data["maxMonthlySession"]; ok {
		upsertOperation("radcheck", "Max-Monthly-Session", value, userinfo.UserName)
	}

	if value, ok := data["maxDailySession"]; ok {
		upsertOperation("radcheck", "Max-Daily-Session", value, userinfo.UserName)
	}

	if value, ok := data["totalData"]; ok {
		upsertOperation("radcheck", "CoovaChilli-Max-Total-Octets", value, userinfo.UserName)

	}

	if value, ok := data["totalDataMonthly"]; ok {
		upsertOperation("radcheck", "CoovaChilli-Max-Total-Octets-Monthly", value, userinfo.UserName)
	}

	if value, ok := data["totalDataDaily"]; ok {
		upsertOperation("radcheck", "CoovaChilli-Max-Total-Octets-Daily", value, userinfo.UserName)
	}

	if value, ok := data["timeOut"]; ok {
		upsertOperationReplay("radreply", "Session-Timeout", value, userinfo.UserName)
	}
	if value, ok := data["maxUpBandWidth"]; ok {
		upsertOperationReplay("radreply", "CoovaChilli-Bandwidth-Max-Up", value, userinfo.UserName)
	}
	if value, ok := data["maxDownBandWidth"]; ok {
		upsertOperationReplay("radreply", "CoovaChilli-Bandwidth-Max-Down", value, userinfo.UserName)
	}
	var tmp CreateResponse
	tmp.UserName = userinfo.UserName
	tmp.Id = userinfo.Id
	tmp.Mail = userinfo.Mail
	return c.JSON(fiber.Map{
		"message": "success",
		"user":    tmp,
	})

}

func UpdateUser(c *fiber.Ctx) error {

	var data map[string]string

	if error := c.BodyParser(&data); error != nil {
		return error
	}

	var checkuser models.UserInfo

	database.DB.Table("userinfo").Where("username = ?", data["username"]).First(&checkuser)

	if checkuser.UserName == data["username"] {

		if value, ok := data["sessionAllowed"]; ok {
			upsertOperation("radcheck", "Max-All-Session", value, checkuser.UserName)
		}

		if value, ok := data["maxMonthlySession"]; ok {
			upsertOperation("radcheck", "Max-Monthly-Session", value, checkuser.UserName)
		}

		if value, ok := data["maxDailySession"]; ok {
			upsertOperation("radcheck", "Max-Daily-Session", value, checkuser.UserName)
		}

		if value, ok := data["totalData"]; ok {
			upsertOperation("radcheck", "CoovaChilli-Max-Total-Octets", value, checkuser.UserName)
		}

		if value, ok := data["totalDataMonthly"]; ok {
			upsertOperation("radcheck", "CoovaChilli-Max-Total-Octets-Monthly", value, checkuser.UserName)
		}

		if value, ok := data["totalDataDaily"]; ok {
			upsertOperation("radcheck", "CoovaChilli-Max-Total-Octets-Daily", value, checkuser.UserName)
		}
		if value, ok := data["timeOut"]; ok {
			upsertOperationReplay("radreply", "Session-Timeout", value, checkuser.UserName)
		}
		if value, ok := data["maxUpBandWidth"]; ok {
			upsertOperationReplay("radreply", "CoovaChilli-Bandwidth-Max-Up", value, checkuser.UserName)
		}
		if value, ok := data["maxDownBandWidth"]; ok {
			upsertOperationReplay("radreply", "CoovaChilli-Bandwidth-Max-Down", value, checkuser.UserName)
		}
		if value, ok := data["password"]; ok {
			updatePassword(checkuser.UserName, value)
		}

		return c.JSON(fiber.Map{
			"message": "success",
		})
	} else {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "username not exist",
		})
	}

}
func updatePassword(username, password string) {
	var str string
	var shuf string
	str = php2go.StrShuffle("abcdefghijklmnopqrstuvwxyz0123456789")
	shuf = php2go.Substr(str, 0, 8)
	fmt.Println(shuf)
	salt := fmt.Sprintf("%s%s%s", []byte("$5$"), []byte(shuf), []byte("$"))
	fmt.Println(salt)

	encoded, err := crypt.Crypt(password, salt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(encoded)
	upsertOperation("radcheck", "Crypt-Password", encoded, username)

}

func upsertOperation(tableName, column, value, userName string) {
	var field models.RedCheck
	field.UserName = userName
	field.Attribute = column
	field.Operation = ":="
	field.Value = value

	var listA []models.RedCheck
	var count int64
	database.DB.Table(tableName).Where(map[string]interface{}{"username": userName,
		"value": value, "attribute": column}).Find(&listA).Count(&count)

	if count == 0 {
		if database.DB.Table(tableName).Where(map[string]interface{}{"username": userName, "attribute": column}).Updates(&field).RowsAffected == 0 {
			database.DB.Table(tableName).Create(&field)
		}
	}

	return
}
func upsertOperationReplay(tableName, column, value, userName string) {
	var field models.RedReplay
	field.UserName = userName
	field.Attribute = column
	field.Operation = "="
	field.Value = value

	var listA []models.RedReplay
	var count int64
	database.DB.Table(tableName).Where(map[string]interface{}{"username": userName,
		"value": value, "attribute": column}).Find(&listA).Count(&count)

	if count == 0 {
		if database.DB.Table(tableName).Where(map[string]interface{}{"username": userName, "attribute": column}).Updates(&field).RowsAffected == 0 {
			database.DB.Table(tableName).Create(&field)
		}
	}

	return
}
func GetAllUser(c *fiber.Ctx) error {

	var userInfoList []models.UserInfo
	database.DB.Table("userinfo").Find(&userInfoList)

	return c.JSON(userInfoList)

}
