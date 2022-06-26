package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	waitTime := 3
	counter := 0

	flag.Parse()
	args := flag.Args()
	isOK := make(chan struct{})
Loop:
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		err := cmd.Start()
		if err != nil {
			log.Println("Failed Start Command.")
			continue Loop
		}
		counter += 1
		log.Printf("Start Command(%d)\n", counter)
		go func(c int) {
			if err = cmd.Wait(); err != nil {
				log.Printf("Failed Run Command(%d): %v\n", c, err)
				return
			}
			isOK <- struct{}{}
		}(counter)
		for {
			select {
			case <-isOK:
				cancel()
				fmt.Println("Succeeded Command!")
				break Loop
			case <-ctx.Done():
				cancel()
				log.Println("Failed: ", ctx.Err())
				continue Loop
			}
		}
	}
}
