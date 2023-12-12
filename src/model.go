package main

import "time"

type Task struct {
	User        string    `json:"user" db:"user"`
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ConditionId int       `json:"condition_id" db:"condition_id"`
	Difficulty  int       `json:"difficulty" db:"difficulty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
}
