package exercises

import (
	"fmt"
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
	terminated := doWork(done, nil) //here we would do nothing for that nil value and be a deadlock

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")

}
