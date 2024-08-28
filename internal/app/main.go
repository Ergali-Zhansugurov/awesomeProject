package app

import (
	configs "awesomeProject/internal/config"
	"awesomeProject/internal/http"
	"awesomeProject/internal/message_broker/broker"
	"awesomeProject/internal/store/postgres"
	"context"
	"fmt"
	"log"

	lru "github.com/hashicorp/golang-lru"

	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)
	cfg := configs.GetConfig()
	store := postgres.NewDB()
	if err := store.Connect(&cfg.Storage); err != nil {
		panic(err)
	}
	defer store.Close()

	fmt.Println("The bd connected")

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	broker := broker.NewBroker(cfg.MainConfig, cache, "Name")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	fmt.Println("The broker connected")

	srv := http.NewServer(
		ctx,
		http.WithAddress(cfg.Listen.Port),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGraceFulTarmination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}
