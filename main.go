package main

import (
	"log"

	//"github.com/gofiber/template/html/v2"

	/* 	jwtware "github.com/gofiber/contrib/jwt" */
	"github.com/gofiber/fiber/v2"
	"github.com/greybluesea/jwt-auth-gofiber/database"
	routes "github.com/greybluesea/jwt-auth-gofiber/routes"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()
	//	engine := html.New("./views", ".html")
	app := fiber.New(
	/* 	fiber.Config{
		Views:       engine,
		ViewsLayout: "layout",
	} */)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, welcome to the JWT auth GoFiber api 👋!")
	})

	/* 	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
	})) */

	routes.SetAuthRoutes(app)
	routes.SetUserRoutes(app)
	routes.SetSigninRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
