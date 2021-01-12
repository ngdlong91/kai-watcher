package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Second().Do(func() {
		fmt.Println("Test", time.Now())
		time.Sleep(5 * time.Second)
	})
	scheduler.StartAsync()
	scheduler.StartBlocking()

}
