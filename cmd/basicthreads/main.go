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

type Category struct {
	ID            int
	Name          string
	ParentID      int
	Subcategories []Category
}

func register(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	password := c.FormValue("password")

	response := users.RegisterUser(name, email, phone, password)

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
	return c.JSON(http.StatusOK, products)
}

func get_products_category(c echo.Context) error {
	id := c.Param("id")
	products := database.GetProductsCategory(id)

	return c.JSON(http.StatusOK, products)
}

func get_categories(c echo.Context) error {
	dbcategories := database.GetCategories("")
	categories := make([]Category, len(dbcategories))

	for i, category := range dbcategories {

		categories[i] = Category{
			ID:       category.ID,
			Name:     category.Name,
			ParentID: category.ParentID,
		}

		stringID := fmt.Sprintf("%d", category.ID)
		subcategories := database.GetCategories(stringID)
		for _, subcategory := range subcategories {
			categories[i].Subcategories = append(categories[i].Subcategories, Category{
				ID:       subcategory.ID,
				Name:     subcategory.Name,
				ParentID: subcategory.ParentID,
			})
		}
	}

	return c.JSON(http.StatusOK, categories)
}

func get_product(c echo.Context) error {
	id := c.Param("id")
	product := database.GetProduct(id)
	return c.JSON(http.StatusOK, product)
}

func get_category(c echo.Context) error {
	id := c.Param("id")
	product := database.GetCategoryName(id)
	return c.JSON(http.StatusOK, product)
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
	e.GET("/products/:id", get_products_category)
	e.GET("/product/:id", get_product)
	e.GET("/categories", get_categories)
	e.GET("/categories/:id", get_category)

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
