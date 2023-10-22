package main

import (
	"log"

	/* 	jwtware "github.com/gofiber/contrib/jwt" */
	"github.com/gofiber/fiber/v2"
	"github.com/greybluesea/jwt-auth-gofiber/database"
	routes "github.com/greybluesea/jwt-auth-gofiber/routes"
)

func main() {

	database.ConnectDB()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, welcome to the JWT auth GoFiber server")
	})
	/* 	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
	})) */

	routes.SetAuthRoutes(app)
	routes.SetUserRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
