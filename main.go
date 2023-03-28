package main

import (
	"fmt"
	"sync"
)

var m sync.Mutex
var set = make(map[int]bool, 0)

func printOnce(num int) {
	m.Lock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
	m.Unlock()
}
func add(keys ...int) int {
	sum := 0
	for _, key := range keys {
		sum = sum + key
	}
	return sum
}
func main() {
	//for i := 0; i < 10; i++ {
	//	go printOnce(100)
	//}
	//time.Sleep(time.Second)
	//log.Println("one")

}
