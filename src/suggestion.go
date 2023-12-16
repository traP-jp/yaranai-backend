package main

import (
	"database/sql"

	"time"
)

func Suggest(user string) (tasks []Task, err error) {
	// now := time.Now()
	// connect to database as root@localhost and open deleted_task table
	db, err := sql.Open("mysql", "yaranai@tcp(localhost:3306)/deleted_task")
	if err != nil {
		return
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
	return
}

func getTimeSlotsForClustering(user string) (time_slots_for_clustering []TimeSlotForClustering, err error) {
	// connect to database as root@localhost and open deleted_task table
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/deleted_task")
	if err != nil {
		return
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
	// accumulate all condition ids in unordered set
	condition_ids := make(map[int]bool)
	for _, deleted_task_property := range deleted_task_properties {
		condition_ids[deleted_task_property.ConditionId] = true
	}
	// construct []int from unordered set
	condition_ids_slice := make([]int, 0)
	for condition_id, _ := range condition_ids {
		condition_ids_slice = append(condition_ids_slice, condition_id)
	}
	// for each time slot, count the number of tasks with each condition id
	condition_ids_distribution := make(map[int]map[int]map[int]int)
	for _, deleted_task_property := range deleted_task_properties {
		deleted_time := deleted_task_property.DeletedAtUnix
		deleted_day_of_week := time.Unix(deleted_time, 0).Weekday()
		deleted_hour_of_day := time.Unix(deleted_time, 0).Hour()
		condition_ids_distribution[int(deleted_day_of_week)][int(deleted_hour_of_day)][deleted_task_property.ConditionId]++
	}
	// construct []TimeSlotForClustering
	time_slots := make([]TimeSlotForClustering, 0)
	for day := 0; day < 7; day++ {
		for hour := 0; hour < 24; hour++ {
			time_slots = append(time_slots, TimeSlotForClustering{
				DeletedDayOfWeek:        day,
				DeletedHourOfDay:        hour,
				ConditionIds:            condition_ids_slice,
				ConditionIdDistribution: condition_ids_distribution[day][hour],
			})
		}
	}
	return time_slots, nil
}
