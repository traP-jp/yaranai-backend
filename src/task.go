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
	tasks := []Task{}
	if err := db.Select(&tasks, "SELECT * FROM task WHERE user = ?", userId); err != nil {
		fmt.Println(err)
	}
	res := []TaskRes{}
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

func deleteTaskHandler(c echo.Context) error {
	var payload AuthHeader
	(&echo.DefaultBinder{}).BindHeaders(c, &payload)
	userId := payload.UserId

	taskId := c.Param("id")

	var task Task
	if err := db.Get(&task, "SELECT * FROM task WHERE id = ?", taskId); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get task"})
	}
	if task.User != userId {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "This task is not yours"})
	}

	_, err := db.Exec("DELETE FROM task WHERE id = ?", taskId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete task"})
	}

	// suggestion : add DeletedTask to deleted_task table
	dateOfNow := time.Now()
	deleting_task := DeletedTask{
		User:          task.User,
		Id:            task.Id,
		ConditionId:   task.ConditionId,
		CreatedAt:     task.CreatedAt,
		DueDate:       task.DueDate,
		DeletedAtUnix: dateOfNow.Unix(),
	}
	_, err = db.Exec("INSERT INTO deletd_task (user, id, condition_id, created_at, due_date, deleted_at_unix) VALUES (?,?,?,?,?,?)", deleting_task.User, deleting_task.Id, deleting_task.ConditionId, deleting_task.CreatedAt, deleting_task.DueDate, deleting_task.DeletedAtUnix)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add deleted task"})
	}

	return c.JSON(http.StatusOK, TaskRes{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		ConditionId: task.ConditionId,
		Difficulty:  task.Difficulty,
		DueDate:     task.DueDate.Format("2006-01-02"),
	})
}
