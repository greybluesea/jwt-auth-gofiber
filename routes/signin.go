package routes

import "github.com/gofiber/fiber/v2"

func SetSigninRoutes(app *fiber.App) {

	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("signup", fiber.Map{"Title": "Sign Up"})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{"Title": "Log In"})
	})

}
