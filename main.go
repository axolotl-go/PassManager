package main

import (
	"fmt"

	"github.com/axolotl-go/password-manager/db"
	"github.com/axolotl-go/password-manager/models"
	"github.com/axolotl-go/password-manager/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Dbconnection()
	db.DB.AutoMigrate(models.PasswordEntry{})
	db.DB.AutoMigrate(models.User{})

	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))

	app.Post("/user/", routes.PostUserHandler)
	app.Get("/users/", routes.GetAllUsers)

	app.Listen(":8080")
	fmt.Println("open port...")
}
