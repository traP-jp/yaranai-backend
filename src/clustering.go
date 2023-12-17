package main

import (
	"math"
	"time"

	"github.com/sohlich/go-dbscan"
)

func (this TimeSlotForClustering) Distance(c interface{}) float64 {
	this_time := this.DeletedDayOfWeek*24 + this.DeletedHourOfDay
	c_time := c.(TimeSlotForClustering).DeletedDayOfWeek*24 + c.(TimeSlotForClustering).DeletedHourOfDay
	diff := math.Abs(float64(this_time - c_time))
	distance := math.Min(diff, 168-diff)
	return distance
}

func (this TimeSlotForClustering) GetID() string {
	return this.UUID
}

func clusterize(time_slots_for_clustering []TimeSlotForClustering) [][]TimeSlotForClustering {
	// clustering
	clusters_list := []dbscan.Clusterable{}
	for _, time_slot_for_clustering := range time_slots_for_clustering {
		clusters_list = append(clusters_list, dbscan.Clusterable(time_slot_for_clustering))
	}
	clusters := dbscan.Clusterize(clusters_list, 24, 2)
	// convert clusters to [][]TimeSlotForClustering
	clusterized_time_slots_for_clustering := make([][]TimeSlotForClustering, 0)
	for _, cluster := range clusters {
		clusterized_time_slots_for_clustering = append(clusterized_time_slots_for_clustering, []TimeSlotForClustering{})
		for _, clusterable := range cluster {
			clusterized_time_slots_for_clustering[len(clusterized_time_slots_for_clustering)-1] = append(clusterized_time_slots_for_clustering[len(clusterized_time_slots_for_clustering)-1], clusterable.(TimeSlotForClustering))
		}
	}
	return clusterized_time_slots_for_clustering
}

func getProbabilityOfNow(clusters [][]TimeSlotForClustering, condition_ids []int) map[int]float64 {
	now := time.Now()
	day_of_week := now.Weekday()
	hour_of_day := now.Hour()
	// find cluster which contains now
	var cluster []TimeSlotForClustering = nil
	for _, c := range clusters {
		for _, time_slot_for_clustering := range c {
			if time_slot_for_clustering.DeletedDayOfWeek == int(day_of_week) && time_slot_for_clustering.DeletedHourOfDay == int(hour_of_day) {
				cluster = c
				break
			}
		}
	}
	probabilities := make(map[int]float64)
	if cluster == nil {
		for _, condition_id := range condition_ids {
			probabilities[condition_id] = 1.0 / float64(len(condition_ids))
		}
	} else {
		count_of_cases := 0
		for _, time_slot_for_clustering := range cluster {
			for condition_id, count := range time_slot_for_clustering.ConditionIdDistribution {
				probabilities[condition_id] += float64(count)
				count_of_cases += count
			}
		}
		for condition_id := range probabilities {
			probabilities[condition_id] /= float64(count_of_cases)
		}
	}
	return probabilities
}

func getTimeSlotsForClustering(user string, deleted_tasks []DeletedTask) (time_slots_for_clustering []TimeSlotForClustering, err error) {

	// accumulate all condition ids in unordered set
	condition_ids := make(map[int]bool)
	for _, deleted_task_property := range deleted_tasks {
		condition_ids[deleted_task_property.ConditionId] = true
	}
	// construct []int from unordered set
	condition_ids_slice := make([]int, 0)
	for condition_id := range condition_ids {
		condition_ids_slice = append(condition_ids_slice, condition_id)
	}
	// for each time slot, count the number of tasks with each condition id
	condition_ids_distribution := make(map[int]map[int]map[int]int)
	for _, deleted_task_property := range deleted_tasks {
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
