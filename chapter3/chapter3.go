package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	errG, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/geek_time", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("geek_time"))
	})

	server := http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:5001",
	}

	errG.Go(func() error {
		return server.ListenAndServe()
	})

	errG.Go(func() error {
		timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		quitSignal := make(chan os.Signal)
		signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case signalQuit := <-quitSignal:
			//
			_ = server.Shutdown(timeoutCtx)
			return errors.Errorf("quit signal: %v", signalQuit)
		}
	})
	signalQuit := errG.Wait()
	fmt.Printf("geektime errgroup quit: %+v\n", signalQuit)
}
