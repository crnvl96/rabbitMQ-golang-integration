package main

import "github.com/crnvl96/rabbitmq-golang/pkg/rabbitmq"

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbitmq.Publish(ch, "Hello!", "amq.direct")
}
