package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func TestSuggest(t *testing.T) {
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
	task, err := suggest("ramdos", conf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(task) != 0 {
		t.Errorf("unexpected task: %v", task)
	}
	for _, task := range task {
		fmt.Println(task)
	}
}
