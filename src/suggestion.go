package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
)

func suggest(user string) ([]Task, error) {
	db, err := sql.Open("mysql", "yaranai@localhost/deleted_task")
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
	tasks_sorted_by_preference, err := suggestInternal(user, tasks)
	if err != nil {
		// if error occurs, stderr error message and return tasks_not_sorted
		fmt.Fprintln(os.Stderr, err)
		return tasks, nil
	} else {
		return tasks_sorted_by_preference, nil
	}
}

func suggestInternal(user string, tasks []Task) ([]Task, error) {
	// now := time.Now()
	// connect to database as root@localhost and open deleted_task table
	db, err := sql.Open("mysql", "yaranai@localhost/deleted_task")
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
	var deleted_task_properties []DeletedTask
	for rows.Next() {
		var deleted_task_property DeletedTask
		err = rows.Scan(
			&deleted_task_property.User,
			&deleted_task_property.Id,
			&deleted_task_property.ConditionId,
			&deleted_task_property.CreatedAt,
			&deleted_task_property.DueDate,
			&deleted_task_property.DeletedAtUnix,
		)
		if err != nil {
			return nil, err
		}
		deleted_task_properties = append(deleted_task_properties, deleted_task_property)
	}
	// get time slots for clustering
	time_slots_for_clustering, err := getTimeSlotsForClustering(user)
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
