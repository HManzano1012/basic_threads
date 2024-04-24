package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"basicthreads/internal/database"
	"basicthreads/internal/users"
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

	response := users.RegisterUser(name, email, phone)

	return c.JSON(http.StatusOK, response)
}

func login(c echo.Context) error {
	username := c.FormValue("email")
	password := c.FormValue("password")

	response := users.LoginUser(username, password)

	return c.JSON(http.StatusOK, response)
}

func get_products(c echo.Context) error {
	products := database.GetProducts()

	fmt.Println(products)

	return c.JSON(http.StatusOK, products)
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
	e.GET("/products", get_products)

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

	e.Logger.Fatal(e.Start(":1323"))
}
