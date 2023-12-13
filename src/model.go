package main

type AuthHeader struct {
	UserId string `header:"X-Showcase-User"`
}
type Task struct {
	User        string `json:"user" db:"user"`
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ConditionId int    `json:"condition" db:"condition_id"`
	Difficulty  int    `json:"difficulty" db:"difficulty"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	DueDate     string `json:"dueDate" db:"dueDate"`
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
