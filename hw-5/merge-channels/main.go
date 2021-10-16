package main

import "sync"

func merge(cs ...<-chan int) <-chan int {
	ch := make(chan int)
	wg := new(sync.WaitGroup)
	wg.Add(len(cs))

	for _, c := range cs {

		localC := c
		go func() {
			defer wg.Done()

			for in := range localC {
				ch <- in
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
