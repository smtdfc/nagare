package helpers

import "sync"

func MergeChannel[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
