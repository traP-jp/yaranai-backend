package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var (
	db *sqlx.DB
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("NS_MARIADB_USER") == "" {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(os.Getenv("NS_MARIADB_USER"))
	fmt.Println("aa")
	conf := mysql.Config{
		User:                 os.Getenv("NS_MARIADB_USER"),
		Passwd:               os.Getenv("NS_MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("NS_MARIADB_HOSTNAME") + ":" + os.Getenv("NS_MARIADB_PORT"),
		DBName:               os.Getenv("NS_MARIADB_DATABASE"),
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	_db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conntected")
	db = _db

	e := echo.New()
	e.GET("/ping", pingHandler)
	e.GET("/tasktest", getTaskTestHandler)
	e.GET("/task", getTaskHandler)

	e.GET("/condition", getConditionHandler)
	
	e.Start(":8080")
}
func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
func getTaskTestHandler(c echo.Context) error {
	var tasks []Task
	if err := db.Select(&tasks, "SELECT * FROM task"); err != nil {
		fmt.Println(err)
	}
	return c.JSON(http.StatusOK, tasks)
}
