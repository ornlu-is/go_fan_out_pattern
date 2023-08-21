package main

import "fmt"

func someNumberStrings() <-chan string {
	ch := make(chan string)
	numberStrings := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	go func() {
		for _, numberString := range numberStrings {
			ch <- numberString
		}

		close(ch)
		return
	}()

	return ch
}

func demultiplexer(worker string, ch <-chan string) chan struct{} {
	stop := make(chan struct{})

	go func() {
		defer close(stop)
		for v := range ch {
			fmt.Println(worker, v)
		}
	}()

	return stop
}

func main() {
	originalStream := someNumberStrings()

	demuxedStream1 := demultiplexer("demux1", originalStream)
	demuxedStream2 := demultiplexer("demux2", originalStream)

	for {
		if demuxedStream1 == nil && demuxedStream2 == nil {
			break
		}

		select {
		case _, ok := <-demuxedStream1:
			if !ok {
				demuxedStream1 = nil
			}
		case _, ok := <-demuxedStream2:
			if !ok {
				demuxedStream2 = nil
			}
		}
	}

	fmt.Println("bye")
}
