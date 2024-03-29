package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func connect() *sql.DB {
	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	dbhost := os.Getenv("DBHOST")
	dbname := os.Getenv("DBNAME")

	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp("+dbhost+")/"+dbname)
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
		"SELECT mail, password FROM users WHERE username = ? AND active = 1 ",
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
