package controllers

import (
	"fmt"
	"mysqlapi/database"
	"mysqlapi/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

type LoginResponse struct {
	UserName string `json:"userName"`
	Id       uint   `json:"id"`
	Email    string `json:"email`
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if error := c.BodyParser(&data); error != nil {
		return error
	}

	if key, ok := data["key"]; !ok {
		fmt.Println(key)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Registration Key Required",
		})
	}
	if data["key"] != "XYZ" {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Registration Key Invalid",
		})
	}
	var checkuser models.User

	database.DB.Where("email = ?", data["email"]).First(&checkuser)

	if checkuser.Email == data["email"] {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "email already exist",
		})
	}

	database.DB.Where("username = ?", data["username"]).First(&checkuser)

	if checkuser.UserName == data["username"] {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "username already exist",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		UserName: data["username"],
		Password: string(password),
	}

	database.DB.Create(&user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	var tmp LoginResponse
	tmp.UserName = user.UserName
	tmp.Id = user.Id
	tmp.Email = user.Email

	return c.JSON(fiber.Map{
		"message": "success",
		"user":    tmp,
		"token":   token,
	})

}

func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("userName = ?", data["username"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	var tmp LoginResponse
	tmp.UserName = user.UserName
	tmp.Id = user.Id
	tmp.Email = user.Email

	return c.JSON(fiber.Map{
		"message": "success",
		"token":   token,
		"user":    tmp,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	var tmp LoginResponse
	tmp.UserName = user.UserName
	tmp.Id = user.Id
	tmp.Email = user.Email
	return c.JSON(user)

}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})

}
func ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])

		}
		return []byte(SecretKey), nil
	})
}
