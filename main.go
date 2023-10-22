package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/greybluesea/jwt-auth-gofiber/database"
	"github.com/greybluesea/jwt-auth-gofiber/routes"
)

func main() {

	database.ConnectDB()

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
