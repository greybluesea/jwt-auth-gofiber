package routes

import (
	//"fmt"
	"log"
	"os"
	"time"

	//	"github.com/gofiber/contrib/jwt"
	//"crypto/rand"
	//"crypto/rsa"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/greybluesea/jwt-auth-gofiber/database"
	"github.com/greybluesea/jwt-auth-gofiber/models"
	"golang.org/x/crypto/bcrypt"
)

/* var (
	// Obviously, this is just a test example. Do not do this in production.
	// In production, you would have the private key and public key pair generated
	// in advance. NEVER add a private key to any GitHub repo.
	privateKey *rsa.PrivateKey
) */

func SetAuthRoutes(app *fiber.App) {

	/* 	// Just as a demo, generate a new private/public key pair on each run. See note above.
	   	rng := rand.Reader
	   	var err error
	   	privateKey, err = rsa.GenerateKey(rng, 2048)
	   	if err != nil {
	   		log.Fatalf("rsa.GenerateKey: %v", err)
	   	} */

	groupAuth := app.Group("/auth")

	groupAuth.Post("/signup", func(c *fiber.Ctx) error {
		signup := new(models.SignupRequest)

		/* signup.Email = c.FormValue("Email")
		signup.Name = c.FormValue("Name")
		signup.Password = c.FormValue("Password") */

		if err := c.BodyParser(&signup); err != nil {
			return err
		}

		if signup.Name == "" || signup.Email == "" || signup.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid sign-up credentials")
		}

		// save this info in the database
		hash, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Name:           signup.Name,
			Email:          signup.Email,
			HashedPassword: string(hash),
		}

		result := database.DB.Create(&user)
		if result.Error != nil {
			return result.Error
		}

		token, err := createJWTTokenSTr(&user)
		if err != nil {
			return err
		}

		c.Cookie(&fiber.Cookie{
			Name:  "jwt",
			Value: token,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{ /* "status": "success", "message": "Sign-up success", */ "token": token})
	})

	groupAuth.Post("/login", func(c *fiber.Ctx) error {
		login := models.LoginRequest{}

		login.Email = c.FormValue("Email")
		login.Password = c.FormValue("Password")

		/* if err := c.BodyParser(&login); err != nil {
			return err
		} */

		if login.Email == "" || login.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid login credentials")
		}

		user := models.User{}
		database.DB.Find(&user, "Email = ?", login.Email)
		if user.ID == 0 {
			return c.Status(400).JSON("this email is not registered")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password)); err != nil {
			return err
		}

		token, err := createJWTTokenSTr(&user)
		if err != nil {
			return err
		}

		c.Cookie(&fiber.Cookie{
			Name:  "jwt",
			Value: token,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{ /* "status": "success", "message": "Log-in success", */ "token": token})
	})

}

func createJWTTokenSTr(user *models.User) (string, error) {

	claims := jwt.MapClaims{
		"name": user.Name,
		//"admin": true,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new JWT token with the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the JWT token using a secret key and get the token string
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))

	/* // Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	tokenStr, err := token.SignedString(privateKey)

	*/
	// If there's an error while signing the token, return an error
	if err != nil {
		log.Fatal("token.SignedString: %w", err)
	}

	// Return the JWT token string, expiration time, and no error
	return tokenStr, nil
}
