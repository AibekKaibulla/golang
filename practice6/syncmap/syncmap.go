package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 101; i++ {
		wg.Add(1)
		go func (key int)  {
			defer wg.Done()
			sm.Store("key", key)
		}(i)
	}

	wg.Wait()
	val, _ := sm.Load("key")
	fmt.Printf("Value: %v\n", val)
}
