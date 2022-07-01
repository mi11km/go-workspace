package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os/exec"
	"time"
)

func main() {
	waitTime := 10
	counter := 0

	flag.Parse()
	args := flag.Args()
	isOK := make(chan struct{})
Loop:
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(waitTime))
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Start()
		if err != nil {
			log.Println("[RETRY] Failed Start Command.")
			break Loop
		}

		counter += 1
		log.Printf("[RETRY] Start Command(%d)\n", counter)

		go func(c int) {
			if err = cmd.Wait(); err != nil {
				log.Printf("[RETRY] Failed Run Command(%d): %v\n", c, err)
				if stderr.Len() != 0 {
					log.Println(stderr.String())
				}
				return
			}
			isOK <- struct{}{}
		}(counter)

		for {
			select {
			case <-isOK:
				cancel()
				log.Println("[RETRY] Succeeded Command!")
				log.Println(stdout.String())
				break Loop
			case <-ctx.Done():
				cancel()
				log.Println("[RETRY] Failed: ", ctx.Err())
				continue Loop
			default:
				if stdout.Len() != 0 {
					log.Println(stdout.String())
				}
			}
		}
	}
}
