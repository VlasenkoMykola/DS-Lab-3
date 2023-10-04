package main

import (
	"fmt"
	"sync"
)

type counter struct {
	sync.Mutex
	count int
}

func (c *counter) Increment() {
	c.Lock()
	defer c.Unlock()
	c.count++
}

func (c *counter) Decrement() {
	c.Lock()
	defer c.Unlock()
	c.count--
}

func (c *counter) getCounter() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}

func awakeAndEat(ch chan int) {
	close(ch)
	fmt.Printf("Honey Pot is full %v\n", len(ch))
	fmt.Println("Winnie the Pooh is awake")
	for i := range ch {
		fmt.Printf("Winnie the Pooh is eating honey %v\n", i)
	}
}

var QUIT = new(counter)
var wg sync.WaitGroup

func fill(ch chan int) {
	defer wg.Done()
	if QUIT.getCounter() == 1 {
		return
	}
	select {
	case ch <- len(ch):
		fmt.Printf("Filling the honey pot %v\n", len(ch))
	default:
		QUIT.Increment()
		awakeAndEat(ch)
	}
}

func main() {
	CAPACITY_OF_HONEY_POT := 100
	AMOUNT_OF_BEES := 25

	ch := make(chan int, CAPACITY_OF_HONEY_POT)

	for QUIT.getCounter() == 0 {
		for i := 0; i < AMOUNT_OF_BEES; i++ {
			if QUIT.getCounter() == 1 {
				break
			}
			fmt.Printf("Bee is heading to jug %v\n", i)
			wg.Add(1)
			go fill(ch)
		}
	}
	wg.Wait()
	fmt.Println("Winnie the Pooh ate everything...")
}
