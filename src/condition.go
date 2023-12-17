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

	conditions := []ConditionWithoutUser{}
	if err := db.Select(&conditions, "SELECT `condition_id`, `condition` FROM `condition` WHERE `user`=?", userId); err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, conditions)
}

func postConditionHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	//JSONリクエストボディのBind
	conditionreq := &ConditionRequestBody{}
	err := c.Bind(conditionreq)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	//空白(0文字)状況登録の制限
	if len(conditionreq.Name) == 0 {
		return c.JSON(http.StatusBadRequest, "Name cannot be empty")
	}

	_, err = db.Exec("INSERT INTO `condition` (`user`,`condition`) VALUES(?,?)", userId, conditionreq.Name)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add condition"})
	}

	//追加した状況のcondition_idを取得
	var condition int
	err = db.Get(&condition, "SELECT `condition_id` FROM `condition` WHERE `condition` = ? ORDER BY `condition_id` DESC", conditionreq.Name)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get condition ID"})
	}

	conditionstr := strconv.Itoa(condition)
	return c.String(http.StatusOK, conditionstr)

}

func putConditionHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	//JSONリクエストボディのBind
	conditionreq := &ConditionRequestBody{}
	err := c.Bind(conditionreq)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	//空白(0文字)状況登録の制限
	if len(conditionreq.Name) == 0 {
		return c.JSON(http.StatusBadRequest, "Name cannot be empty")
	}

	//idのint変換
	putidstr := c.Param("id")
	putid, err := strconv.Atoi(putidstr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	//変更対象conditionの取得
	var condition Condition
	err = db.Get(&condition, "SELECT * FROM `condition` WHERE `condition_id` =?", putid)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	//他ユーザーによる変更を制限
	if condition.User != userId {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "This condition is not yours"})
	}

	//UPDATEの実行
	_, err = db.Exec("UPDATE `condition` set `condition`=? WHERE `condition_id` =?", conditionreq.Name, putid)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to rename condition"})
	}

	return c.JSON(http.StatusOK, ConditionWithoutUser{Id: condition.Id, Name: conditionreq.Name})

}

func deleteConditionHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	//idのint変換
	deleteidstr := c.Param("id")
	deleteid, err := strconv.Atoi(deleteidstr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var condition Condition
	err = db.Get(&condition, "SELECT * FROM `condition` WHERE `condition_id` =?", deleteid)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Failed to find condition"})
	}

	//他ユーザーによる削除を制限
	if condition.User != userId {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "This condition is not yours"})
	}

	//消去実行
	_, err = db.Exec("DELETE FROM `condition` WHERE `condition_id` =?", deleteid)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete condition"})
	}

	return c.JSON(http.StatusOK, ConditionWithoutUser{Id: condition.Id, Name: condition.Name})

}
