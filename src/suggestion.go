package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func suggest(user string, conf mysql.Config) ([]Task, error) {
	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// get all tasks with user
	rows, err := db.Query("SELECT * FROM task WHERE user = ?", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.User,
			&task.Id,
			&task.Title,
			&task.Description,
			&task.ConditionId,
			&task.Difficulty,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DueDate,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	tasks_sorted_by_preference, err := suggestInternal(user, conf, tasks)
	if err != nil {
		// if error occurs, stderr error message and return tasks_not_sorted
		fmt.Fprintln(os.Stderr, err)
		return tasks, err
	} else {
		return tasks_sorted_by_preference, nil
	}
}

func suggestInternal(user string, conf mysql.Config, tasks []Task) ([]Task, error) {
	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// select all deleted_task with user
	rows, err := db.Query("SELECT * FROM deleted_task WHERE user = ?", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deleted_tasks []DeletedTask
	for rows.Next() {
		var deleted_task DeletedTask
		err = rows.Scan(
			&deleted_task.User,
			&deleted_task.Id,
			&deleted_task.ConditionId,
			&deleted_task.CreatedAt,
			&deleted_task.DueDate,
			&deleted_task.DeletedAtUnix,
		)
		if err != nil {
			return nil, err
		}
		deleted_tasks = append(deleted_tasks, deleted_task)
	}
	// get time slots for clustering
	time_slots_for_clustering, err := getTimeSlotsForClustering(user, deleted_tasks)
	if err != nil {
		return nil, err
	}
	// clustering
	clusters := clusterize(time_slots_for_clustering)
	// get probabilities
	// get all condition ids
	condition_ids_unique := make(map[int]bool)
	for _, task := range tasks {
		condition_ids_unique[task.ConditionId] = true
	}
	condition_ids := make([]int, 0)
	for condition_id := range condition_ids_unique {
		condition_ids = append(condition_ids, condition_id)
	}
	probabilities := getProbabilityOfNow(clusters, condition_ids)
	// set task preferences
	type SortStruct struct {
		task       Task
		preference float64
	}
	sort_structs := make([]SortStruct, 0)
	for _, task := range tasks {
		probability_of_condition := probabilities[task.ConditionId]
		preference := taskPreference(task, probability_of_condition)
		sort_struct := SortStruct{
			task:       task,
			preference: preference,
		}
		sort_structs = append(sort_structs, sort_struct)
	}
	// sort by preference
	sort.Slice(sort_structs, func(i, j int) bool {
		return sort_structs[i].preference > sort_structs[j].preference
	})
	// get sorted tasks
	tasks_sorted_by_preference := make([]Task, 0)
	for _, sort_struct := range sort_structs {
		tasks_sorted_by_preference = append(tasks_sorted_by_preference, sort_struct.task)
	}
	return tasks_sorted_by_preference, nil
}
