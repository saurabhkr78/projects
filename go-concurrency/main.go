package main

import (
	"fmt"
	"sync"
)

func Sub(wg *sync.WaitGroup) {
	defer wg.Done()
	a := 10
	b := 20
	diff := a - b
	fmt.Printf("The difference of %d and %d is %d\n", a, b, diff)
}
func sayHello(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello from goroutine")
}
func Sum(wg *sync.WaitGroup) {
	defer wg.Done()
	a := 10
	b := 20
	sum := a + b
	fmt.Printf("The sum of %d and %d is %d\n", a, b, sum)
}
func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go sayHello(&wg)
	go Sum(&wg)
	go Sub(&wg)
	wg.Wait()
	fmt.Println("Main finished")
}
