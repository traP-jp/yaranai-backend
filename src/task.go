package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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
			DueDate:     v.DueDate.Format("2006-01-02"),
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

	dateOfNow := time.Now()
	res, err := db.Exec("INSERT INTO task (user, title, description, condition_id, difficulty, created_at, updated_at, dueDate) VALUES (?,?,?,?,?,?,?,?)", userId, newTask.Title, newTask.Description, newTask.ConditionId, newTask.Difficulty, dateOfNow, dateOfNow, newTask.DueDate)

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
	}

	return c.JSON(http.StatusCreated, addedTask)
}

func putTaskHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	var imformationOfTask TaskWithoutId
	if err := c.Bind(&imformationOfTask); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	var existingTask Task
	err = db.Get(&existingTask, "SELECT * FROM task WHERE id = ? AND user = ?", taskId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
		} else {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve task"})
		}
	}

	dateOfNow := time.Now().Format("2006-01-02")
	_, err = db.Exec("UPDATE task SET title=?, description=?, condition_id=?, difficulty=?, updated_at=?, dueDate=? WHERE id=? AND user=?", imformationOfTask.Title, imformationOfTask.Description, imformationOfTask.ConditionId, imformationOfTask.Difficulty, dateOfNow, imformationOfTask.DueDate, taskId, userId)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update task"})
	}

	imformationOfTaskRes := TaskRes{
		Id:          taskId,
		Title:       imformationOfTask.Title,
		Description: imformationOfTask.Description,
		ConditionId: imformationOfTask.ConditionId,
		Difficulty:  imformationOfTask.Difficulty,
		DueDate:     imformationOfTask.DueDate,
	}

	return c.JSON(http.StatusOK, imformationOfTaskRes)
}
