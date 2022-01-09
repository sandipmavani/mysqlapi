package routes

import (
	"mysqlapi/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Get("/api/logout", controllers.Logout)
	//user operation
	app.Post("/api/user/add", controllers.CreateUser)
	app.Post("/api/user/edit", controllers.UpdateUser)

	app.Post("/api/user/list", controllers.GetAllUser)
}
