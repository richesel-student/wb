package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// прерывание по контексту
func contextFunc(ctx context.Context) {
	fmt.Printf("Работа началась")
	time.Sleep(1 * time.Second)
	<-ctx.Done()
	fmt.Println("Работа прервана :", ctx.Err())

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		fmt.Println("получен сигнал, завершаем ...")
		cancel()
	}()
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	contextFunc(ctx)

}
