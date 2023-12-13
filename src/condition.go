package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getConditionHandler(c echo.Context) error {
	var conditions []Condition
	if err := db.Select(&conditions, "SELECT * FROM `condition`"); err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, conditions)
}

func postConditionHandler(c echo.Context) error {
	conditionreq := &ConditionRequestBody{}
	err := c.Bind(conditionreq)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if len(conditionreq.Name) == 0 {
		return c.JSON(http.StatusBadRequest, "Name cannot be empty")
	}

	_, err = db.Exec("INSERT INTO `condition` (`condition`) VALUES(?)", conditionreq.Name)

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var condition Condition
	err = db.Get(&condition.Id, "SELECT `condition_id` FROM `condition` WHERE `condition` = ? ORDER BY `condition_id` DESC", conditionreq.Name)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, condition)
}
