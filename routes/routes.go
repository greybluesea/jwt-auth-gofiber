package routes

import (
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

		token, exp, err := createJWTToken(&user)
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

func createJWTToken(user *models.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
