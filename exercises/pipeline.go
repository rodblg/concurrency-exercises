package exercises

import (
	"fmt"
)

//a pipeline can be used to process streams or batches of data. You separate concerns in stages. Convert data, stages just like in mongo

var multiply = func(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))

	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}

	return multipliedValues
}

var add = func(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {

		addedValues[i] = v + additive
	}

	return addedValues
}

// You can make a pipeline with both stages
/*

ints := []int{1,2,3,4}
for _, v := range add(multiply(ints,2), 1){
	fmt.Println(v)
}

each stage is taking a slice of data and returning a slice of data. These
stages are performing what we call batch processing.
stream processing. This means that the stage
receives and emits one element at a time

*/

// you can also change the batch processing for stream processing, and applying the stage for an item at a time
var multiplyStream = func(value, multiplier int) int {
	return value * multiplier
}

var addStream = func(value, additive int) int {
	return value + additive
}

/*

ints := []int{1,2,3,4}
for _, v := range ints {

		fmt.Println(addStream(multiplyStream(v,2), 1))

}
*/

// BEST PRACTICES TO CREATE PIPELINES
func StreamCh() {

	//This function converts a discrete set of values into a stream of data on a channel
	var generator = func(done <-chan struct{}, integers ...int) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	var multiplyCh = func(done <-chan struct{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			defer close(multipliedStream)

			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	var addCh = func(
		done <-chan struct{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)

		go func() {
			defer close(addedStream)

			for i := range intStream {
				select {
				case <-done:
					return

				case addedStream <- i + additive:
				}
			}
		}()

		return addedStream
	}

	done := make(chan struct{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5)

	pipeline := multiplyCh(done, addCh(done, multiplyCh(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}

}
func RepeatTakeCh() {
	var repeat = func(
		done <-chan struct{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)

			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()

		return valueStream

	}

	var take = func(
		done <-chan struct{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return

				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	done := make(chan struct{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}

}

var repeatFn = func(
	done <-chan struct{},
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return

			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

//fan-out fan-in

/*
Fan-out is a term to describe the process of starting multiple goroutines to handle
input from the pipeline, and fan-in(multiplexor) is a term to describe the process of combining
multiple results into one channel. Multiplex

reuse a single stage of our pipeline on multiple goroutines in an attempt to parallelize pulls
from an upstream stage
*/
