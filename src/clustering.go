package main

import (
	"math"

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
