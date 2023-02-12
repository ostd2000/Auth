package routes

import "github.com/gofiber/fiber/v2"

func AuthRouter(router fiber.Router) {
	router.Post("/signup")
	router.Post("/login")
	router.Get("/logout")
}