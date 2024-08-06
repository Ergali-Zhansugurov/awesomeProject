package app

import (
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

	dbURL := "postgres://postgres:Fencing.666@localhost:5432/postgres"
	store := postgres.NewDB()
	if err := store.Connect(dbURL); err != nil {
		panic(err)
	}
	defer store.Close()

	fmt.Println("The bd connected")

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	fmt.Println("cache created")

	brokers := []string{"amqp:guest:guest@localhost:5672/"}
	broker := broker.NewBroker(brokers, cache, "Name")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	fmt.Println("The broker connected")

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8080"),
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
