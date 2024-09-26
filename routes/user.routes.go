package routes

import (
	"fmt"

	"github.com/axolotl-go/password-manager/db"
	"github.com/axolotl-go/password-manager/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	// Usar una nueva conexión para esta consulta
	tx := db.DB.Begin()
	defer tx.Rollback()

	if err := tx.Preload("Passwords").Find(&users).Error; err != nil {
		fmt.Println("Error al obtener los usuarios:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error al obtener los usuarios",
			"details": err.Error(),
		})
	}

	tx.Commit()

	// Opcionalmente, ocultar las contraseñas
	for i := range users {
		users[i].Password = ""
		for j := range users[i].Passwords {
			users[i].Passwords[j].Password = "********"
		}
	}

	return c.JSON(users)
}

func PostUserHandler(c *fiber.Ctx) error {
	var userInput struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Passwords []struct {
			Site     string `json:"site"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"passwords"`
	}

	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	user := models.User{
		Email:    userInput.Email,
		Password: userInput.Password, // Considera hashear esta contraseña
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
	}

	for _, pwdEntry := range userInput.Passwords {
		passwordEntry := models.PasswordEntry{
			UserID:   user.ID,
			Site:     pwdEntry.Site,
			Username: pwdEntry.Username,
			Password: pwdEntry.Password,
		}
		if err := db.DB.Create(&passwordEntry).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to create password entry",
				"details": err.Error(),
			})
		}
	}

	// Carga las contraseñas asociadas al usuario
	if err := db.DB.Preload("Passwords").First(&user, user.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to load user with passwords",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
