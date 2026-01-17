package exercises

import (
	"fmt"
	"math/rand"

	"time"
)

//To avoid memory leaks we always have to terminate our go routines

func UnblockingRoutine() {

	//definy anonymous func that has read only channels of interface and string type
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			//we use a for loop select to do something if we receive a signal from done or from strings channel
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	//str := make(chan string)//To send a value we define a channel that would select the s := <-strings case

	terminated := doWork(done, nil) //here we would do nothing for that nil value and be a deadlock
	//str <- "test" but for this to not cause a deadlock a channel must be expecting to receive a value, so that is why it must be after the terminated channel, so the select case of s:=<-strings expects a string
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")

}

func UnblockingRoutineWrite() {
	newRandStream := func(done <-chan struct{}) <-chan int {
		randStream := make(chan int)
		{
			go func() {
				defer fmt.Println("newRandStream closure exited.")
				defer close(randStream)

				for {
					select {
					case randStream <- rand.Int():
					case <-done:
						return
					}

				}
			}()
			return randStream
		}
	}
	done := make(chan struct{})
	randStream := newRandStream(done)

	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d \n", i, <-randStream)
	}

	close(done)
	time.Sleep(1 * time.Second)

}
