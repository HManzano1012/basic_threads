package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"basicthreads/internal/database"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func register(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")

	if len(name) == 0 || len(email) == 0 || len(phone) == 0 {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "Name, email and phone are required",
			"error":   "missing_fields",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	userExists := database.ValidateUserExists(email)
	if userExists {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "User already exists",
			"error":   "user_exists",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	err := database.RegisterUser(name, email, phone)
	if err != nil {
		response := echo.Map{
			"status":  "error",
			"code":    500,
			"message": "Internal server error",
			"error":   "internal_server_error",
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := echo.Map{
		"status":  "success",
		"code":    200,
		"message": "User registered successfully",
	}

	return c.JSON(http.StatusOK, response)
}

func login(c echo.Context) error {
	username := c.FormValue("email")
	password := c.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		response := echo.Map{
			"status":  "error",
			"code":    400,
			"message": "Email and password are required",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	authUser := database.AuthUser(username, password)

	if !authUser {
		response := echo.Map{
			"status":  "error",
			"code":    401,
			"message": "Invalid credentials",
		}

		return c.JSON(http.StatusUnauthorized, response)
	}

	fmt.Println("")

	claims := &jwtCustomClaims{
		username,
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
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"code":   200,
		"token":  t,
	})
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Login route
	e.POST("/login", login)
	e.POST("/register", register)

	// Restricted group
	// r := e.Group("/sms/")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
		ErrorHandler: func(c echo.Context, err error) error {
			response := echo.Map{
				"status":  "error",
				"code":    401,
				"message": "Invalid or expired token",
			}
			return echo.NewHTTPError(http.StatusUnauthorized, response)
		},
	}

	fmt.Println(config)
	// r.Use(echojwt.WithConfig(config))
	// r.POST("send", sendSMS)

	e.Logger.Fatal(e.Start(":1323"))
}
