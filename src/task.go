package main

import (
	"fmt"
	"net/http"
	"time"

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

func postTaskHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	var newTask TaskWithoutId
	if err := c.Bind(&newTask); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	dateOfNow := time.Now().Format("2006-01-02")
	res, err := db.Exec("INSERT INTO task (user, title, description, condition_id, difficulty, created_at, updated_at, dueDate) VALUES (?,?,?,?,?,?)", userId, newTask.Title, newTask.Description, newTask.ConditionId, newTask.Difficulty, dateOfNow, dateOfNow, newTask.DueDate)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add task"})
	}

	taskId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get task ID"})
	}

	addedTask := TaskRes{
		Id:          int(taskId),
		Title:       newTask.Title,
		Description: newTask.Description,
		ConditionId: newTask.ConditionId,
		Difficulty:  newTask.Difficulty,
		DueDate:     newTask.DueDate,
		User:        userId,
	}

	return c.JSON(http.StatusCreated, addedTask)
}
