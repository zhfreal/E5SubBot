package bots

import "testing"

func PerformTasks_test(t *testing.T) {
	for i := 0; i < 10000; i++ {
		PerformTasks()
	}
}
