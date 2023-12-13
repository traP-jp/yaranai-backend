package main

import "time"

type AuthHeader struct {
	UserId string `header:"X-Showcase-User"`
}
type Task struct {
	User        string    `json:"user" db:"user"`
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ConditionId int       `json:"condition" db:"condition_id"`
	Difficulty  int       `json:"difficulty" db:"difficulty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
}
type TaskRes struct {
	User        string    `json:"user" db:"user"`
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ConditionId int       `json:"condition" db:"condition_id"`
	Difficulty  int       `json:"difficulty" db:"difficulty"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
}
type TaskWithoutId struct {
	User        string    `json:"user" db:"user"`
	Description string    `json:"description" db:"description"`
	ConditionId int       `json:"condition" db:"condition_id"`
	Difficulty  int       `json:"difficulty" db:"difficulty"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
}
type Condition struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
