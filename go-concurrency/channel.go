package main

import (
	"fmt"
	"sync"
)

func notification(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	ch <- "notification from goroutine"
}

func printtable(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 5; i++ {
		ch <- fmt.Sprintf("5 x %d = %d\n", i, 5*i)
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan string)
	ch2 := make(chan string)
	go printtable(ch, &wg)
	go notification(ch2, &wg)

	receivedtable := <-ch
	fmt.Println(receivedtable)
	receivedNotification := <-ch2
	fmt.Println(receivedNotification)
}
