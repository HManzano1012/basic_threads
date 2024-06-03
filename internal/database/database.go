package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Category struct {
	ID       int
	Name     string
	ParentID int
}

type Product struct {
	ID          int
	Name        string
	Price       float64
	Description string
	Image       string
}

func connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	dbhost := os.Getenv("DBHOST")
	dbname := os.Getenv("DBNAME")
	dbport := os.Getenv("DBPORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func AuthUser(user string, password string) bool {
	db := connect()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query(
		"SELECT email, password FROM customers WHERE email = ?  ",
		user,
	)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var username string
	var pass string
	if result.Next() {
		result.Scan(&username, &pass)
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	return err == nil
}

func ValidateUserExists(user string) bool {
	db := connect()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
	result, err := db.Query(
		"SELECT email FROM customers WHERE email = ?",
		user,
	)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var username string
	if result.Next() {
		result.Scan(&username)
	}
	return username == user
}

func RegisterUser(name string, email string, phone string) error {
	db := connect()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(
		"INSERT INTO customers (name, email, phone) VALUES (?, ?, ?)",
		name,
		email,
		phone,
	)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func GetProducts() []Product {
	db := connect()
	defer db.Close()

	result, err := db.Query("SELECT product_id as id, name, price, description, img FROM products")
	if err != nil {
		panic(err.Error())
	}

	Products := []Product{}

	for result.Next() {
		var product Product
		err = result.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Image,
		)
		if err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)

	}

	return Products
}

func GetProduct(id string) Product {
	db := connect()
	defer db.Close()

	result, err := db.Query(
		"SELECT product_id as id, name, price, description, img FROM products where product_id = ? ",
		id,
	)
	if err != nil {
		panic(err.Error())
	}

	var product Product
	for result.Next() {
		err = result.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Image,
		)
		if err != nil {
			panic(err.Error())
		}
	}
	return product
}

func GetCategories(id_category string) []Category {
	db := connect()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
	Categories := []Category{}

	if id_category == "" {
		result, err := db.Query(
			"SELECT id,name FROM categories where parent_category_id is null",
		)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var category Category
			err = result.Scan(
				&category.ID,
				&category.Name,
			)
			if err != nil {
				panic(err.Error())
			}

			Categories = append(Categories, category)
		}

	} else {
		result, err := db.Query(
			"SELECT id,name,parent_category_id FROM categories where parent_category_id = ?",
			id_category,
		)
		if err != nil {
			panic(err.Error())
		}

		for result.Next() {
			var category Category
			err = result.Scan(
				&category.ID,
				&category.Name,
				&category.ParentID,
			)
			if err != nil {
				panic(err.Error())
			}

			Categories = append(Categories, category)

		}
	}

	return Categories
}

func GetProductsCategory(id_category string) []Product {
	db := connect()
	defer db.Close()

	result, err := db.Query(
		"SELECT p.product_id as id, p.name, p.price, p.description, p.img FROM products as p inner join categories_product as cp on cp.id_product = p.product_id where cp.id_category = ?",
		id_category,
	)
	if err != nil {
		panic(err.Error())
	}

	Products := []Product{}

	for result.Next() {
		var product Product
		err = result.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Image,
		)
		if err != nil {
			panic(err.Error())
		}

		Products = append(Products, product)

	}

	return Products
}

func GetCategoryName(id string) string {
	db := connect()
	defer db.Close()

	result, err := db.Query(
		"SELECT name FROM categories where id = ?",
		id,
	)
	if err != nil {
		panic(err.Error())
	}
	var name string
	for result.Next() {
		err = result.Scan(&name)
		if err != nil {
			panic(err.Error())
		}
	}
	return name
}
