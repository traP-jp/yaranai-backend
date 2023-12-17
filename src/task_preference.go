package main

func taskPreference(task Task, probability_of_condition float64) float64 {
	return float64(task.Difficulty) * probability_of_condition
}
