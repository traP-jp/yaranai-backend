package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getTaskHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId
	var tasks []Task
	if err := db.Select(&tasks, "SELECT * FROM task WHERE user = ?", userId); err != nil {
		fmt.Println(err)
	}
	var res []TaskRes
	for _, v := range tasks {
		res = append(res, TaskRes{
			Id:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			ConditionId: v.ConditionId,
			Difficulty:  v.Difficulty,
			DueDate:     v.DueDate,
			User:        v.User,
		})
	}
	return c.JSON(http.StatusOK, res)
}
