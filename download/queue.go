package queue

import "fmt"

var queue []int

func Push(val int) {
	queue = append(queue, val)
}

func Pop() {
	if len(queue) > 0 {
		queue = queue[1:]
	}
}
func Front() int {
	if len(queue) == 0 {
		return -1
	}
	return queue[0]
}

func Empty() bool {
	return len(queue) == 0
}