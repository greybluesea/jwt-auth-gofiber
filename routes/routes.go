package routes

import (
	//"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/greybluesea/jwt-auth-gofiber/database"
	"github.com/greybluesea/jwt-auth-gofiber/models"
	"golang.org/x/crypto/bcrypt"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, welcome to the JWT auth GoFiber server")
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		req := new(models.SignupRequest)
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid signup credentials")
		}

		// save this info in the database
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Name:           req.Name,
			Email:          req.Email,
			HashedPassword: string(hash),
		}

		result := database.DB.Create(&user)
		if result.Error != nil {
			return result.Error
		}

		//fmt.Println(result)

		token, exp, err := createJWTTokenSTr(&user)
		if err != nil {
			return err
		}
		// create a jwt token

		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/private", func(c *fiber.Ctx) error {
		return nil
	})

}

func createJWTTokenSTr(user *models.User) (string, int64, error) {
	// Calculate the expiration time for the JWT token (30 minutes from now)
	exp := time.Now().Add(time.Minute * 30).Unix()

	// Create a new JWT token with the HS256 signing method
	token := jwt.New(jwt.SigningMethodHS256)

	// Extract the claims from the token (claims contain the payload of the JWT)
	claims := token.Claims.(jwt.MapClaims)

	// Set the "user_id" claim in the JWT payload to the user's ID
	claims["user_id"] = user.ID

	// Set the "exp" claim in the JWT payload to the expiration time
	claims["exp"] = exp

	// Sign the JWT token using a secret key and get the token string
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))

	// If there's an error while signing the token, return an error
	if err != nil {
		return "", 0, err
	}

	// Return the JWT token string, expiration time, and no error
	return tokenStr, exp, nil
}
