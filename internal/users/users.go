package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"basicthreads/internal/database"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func LoginUser(email, password string) echo.Map {
	if len(email) == 0 || len(password) == 0 {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "Email and password are required",
		}
		return response
	}

	authUser := database.AuthUser(email, password)

	if !authUser {
		response := echo.Map{
			"status":  "error",
			"code":    401,
			"message": "Invalid credentials",
		}

		return response
	}

	claims := &jwtCustomClaims{
		email,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 4)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		response := echo.Map{
			"status":  "error",
			"code":    500,
			"message": "Internal server error",
			"error":   "internal_server_error",
		}

		return response
	}

	response := echo.Map{
		"status": "success",
		"code":   200,
		"token":  t,
	}

	return response
}

func RegisterUser(name, email, phone string) echo.Map {
	if len(name) == 0 || len(email) == 0 || len(phone) == 0 {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "Name, email and phone are required",
			"error":   "missing_fields",
		}

		return response
	}

	userExists := database.ValidateUserExists(email)
	if userExists {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "User already exists",
			"error":   "user_exists",
		}
		return response
	}

	err := database.RegisterUser(name, email, phone)
	if err != nil {
		response := echo.Map{
			"status":  "error",
			"code":    500,
			"message": "Internal server error",
			"error":   "internal_server_error",
		}
		return response
	}

	response := echo.Map{
		"status":  "success",
		"code":    200,
		"message": "User registered successfully",
	}

	return response
}
