package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/PuerkitoBio/goquery"
	//"github.com/opesun/goquery"
)

var (
	WORKERS int    = 1
	url     string = "http://golang-book.ru/chapter-05-control-structures.html"
)

func init() {
	flag.IntVar(&WORKERS, "w", WORKERS, "количество потоков")
	flag.Parse()
}

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

func receive() <-chan []byte {
	rec := make(chan []byte)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			rec <- d.Body
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return rec
}

func par(url string) <-chan string {
	c := make(chan string)
	for i := 0; i < WORKERS; i++ {
		go func() {
			x, err := goquery.NewDocument(url)
			_check(err)
			x.Find("a").Each(func(i int, s *goquery.Selection) {
				link, _ := s.Attr("href")
				c <- link
			})
			time.Sleep(100 * time.Millisecond)
		}()
	}
	fmt.Println("Запущено потоков: ", WORKERS)
	return c
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func send(body string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//body := "aaaaa" /////текст сообщения
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func main() {
	for {
		body_chan := par(url)

		for body_chan != nil {
			body := <-body_chan
			send(body)
		}
		url_chan := receive()
		burl := <-url_chan
		url = string([]byte(burl))

	}
}
