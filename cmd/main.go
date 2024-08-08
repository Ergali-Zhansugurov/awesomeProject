package main

import "awesomeProject/internal/app" //"awesomeProject/internal/app" "manga/internal/app"

func main() {
	/*fmt.Println("Go RabbitMQ test")

	adr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp091.Dial(adr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to Rabbit by adres ", adr)

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"testqueue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(q)

	err = ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Massege Publishet to queue")*/
	app.Run()

}
