package routes

import (
	"fmt"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetUserRoutes(app *fiber.App) {

	userRoutes := app.Group("/user")

	userRoutes.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(os.Getenv("SECRET"))},
		TokenLookup: "cookie:jwt",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Render("unauthorised", fiber.Map{"Title": "Unauthorised, please sign in"}) // Send unauthorized status code and the error
		},
	}))
	/* userRoutes.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    privateKey.Public(),
		},
	})) */

	userRoutes.Get("/me", restricted)

}
func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	return c.Render("home", fiber.Map{"Title": fmt.Sprintln("Welcome ", name)})
	//	return c.SendString("Welcome " + name)

	//return c.SendString("Hello, welcome to the JWT auth GoFiber server")

	//	return c.Render("Authenticated", fiber.Map{"authenticated": true, "path": "user"})
}
