package main

import (
	"fmt"
	"testing"
)

func TestSuggest(t *testing.T) {
	task, err := suggest("ramdos")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(task) != 0 {
		t.Errorf("unexpected task: %v", task)
	}
	for _, task := range task {
		fmt.Println(task)
	}
}
