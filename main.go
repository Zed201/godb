package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"godb/inter"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		inter.ReplCreate()
		sigs <- syscall.SIGINT
	}()
	<-sigs
	fmt.Println("\nSaindo")
}
