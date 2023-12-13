package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func getConditionHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	var conditions []Condition
	if err := db.Select(&conditions, "SELECT * FROM `condition` WHERE `user`=?", userId); err != nil {
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

	var condition int
	err = db.Get(&condition, "SELECT `condition_id` FROM `condition` WHERE `condition` = ? ORDER BY `condition_id` DESC", conditionreq.Name)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	conditionstr := strconv.Itoa(condition)
	return c.String(http.StatusOK, conditionstr)
	
}

func deleteConditionHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId
	
	//idのint変換
	deleteidstr := c.Param("id")
	deleteid, err := strconv.Atoi(deleteidstr)
	if err!= nil {
    return c.String(http.StatusBadRequest, err.Error())
  }

	var condition Condition
	err = db.Get(&condition, "SELECT * FROM `condition` WHERE `condition_id` =?", deleteid)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	//他ユーザーによる削除を制限
	if condition.User != userId {
		return c.String(http.StatusForbidden, "Forbidden")
	}

	//消去実行
	_, err = db.Exec("DELETE FROM `condition` WHERE `condition_id` =?", deleteid)

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, condition)

}
