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
	DueDate     time.Time `json:"dueDate" db:"dueDate"`
}
type TaskRes struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ConditionId int    `json:"condition" db:"condition_id"`
	Difficulty  int    `json:"difficulty" db:"difficulty"`
	DueDate     string `json:"dueDate" db:"dueDate"`
}
type TaskWithoutId struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ConditionId int    `json:"condition" db:"condition_id"`
	Difficulty  int    `json:"difficulty" db:"difficulty"`
	DueDate     string `json:"dueDate" db:"dueDate"`
}
type Condition struct {
	Id int `json:"id" db:"condition_id"`
	User string `json:"user" db:"user"`
	Name string `json:"name,omitempty" db:"condition"`
}
type ConditionRequestBody struct {
	Name string `json:"name,omitempty" db:"condition"`
}
