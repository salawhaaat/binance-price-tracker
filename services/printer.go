package services

import (
	"fmt"
	"time"
)

func StartResultPrinter(wp *WorkerPool) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			wp.mu.Lock()
			fmt.Println(wp.ResultBuffer.String())
			wp.ResultBuffer.Reset()
			fmt.Println("Total requests:", wp.GetRequestsCount())
			fmt.Println("---------------------------\n")
			wp.mu.Unlock()
		}
	}()
}
