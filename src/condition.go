package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)



func getConditionHandler(c echo.Context) error {
	conditions := &[]Condition{}
	if err := db.Select(&conditions, "SELECT * FROM condition"); err != nil {
		fmt.Println(err)
	}
	return c.JSON(http.StatusOK, conditions)
}


