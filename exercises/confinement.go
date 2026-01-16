package exercises

import (
	"bytes"
	"fmt"
	"sync"
)

//you can use data protected by confinement for safe operations
//confinement ensures information is only available from one concurrent process: ad hoc and lexical

// ad hoc, follow a convention to only access data in one point
func AdHocConfinement() {

	data := make([]int, 4)

	//here we only access the values in loopData, so it keeps the synchronization
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)

	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

// lexical confinement, exposes only the correct data for multiple concurrent processes
// we will expose the read or write aspects of a channel that needs them
func LexicalConfinement() {
	chanOwner := func() <-chan int { //confines the write of the channel and then transform the channel into a receive only
		results := make(chan int, 5)
		go func() {
			defer close(results)

			for i := 0; i <= 5; i++ {
				results <- i
			}

		}()

		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}

		fmt.Println("Done receiving!")
	}

	results := chanOwner()

	consumer(results)
}

//channels are concurrent safe so there are no much trouble with lexical confinement

//now we are using a buffer

func ConfinementBuffer() {

	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("example")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}
