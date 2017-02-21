package main

import (
	"flag"
	"github.com/streadway/amqp"
	"log"
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	uri          	= flag.String("uri", "amqp://guest:guest@localhost:15672/", "AMQP URI")
	queue        	= flag.String("queue", "myqueue", "AMQP queue name")
	consumerTag  	= flag.String("consumer-tag", "mycontent-consumer", "AMQP consumer tag")
	driver		= flag.String("driver", "mysql", "SQL Driver")
	dsn		= flag.String("dsn", "root:secret@tcp(localhost:6603)/mydb?charset=utf8", "MySQL DSN")
	insertQuery 	= flag.String("insertQuery", "INSERT mycontent SET message = ? , content = ?", "Insert Query")
)

func init() {
	flag.Parse()
}

type Mycontent struct {
	Message, Content string
}

func main() {
	var messageBody Mycontent

	// Establish MySQL connection
	db, err := sql.Open(*driver, *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if DB is alive!
	err = db.Ping()
	failOnError(err, "Database is not reachable")

	// Prepare insertion
	statement, err := db.Prepare(*insertQuery)
	failOnError(err, "Query could not be prepared")

	// Establish RabbitMQ connection
	connection, err := amqp.Dial(*uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	// Create channel
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	// Register consumer
	messages, err := channel.Consume(
		*queue,
		*consumerTag,
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Create endless channel
	forever := make(chan bool)

	// Process incomming messages
	go func() {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)

			// Unmarshall JSON body
			json.Unmarshal([]byte(message.Body), &messageBody)

			// Insert into DB
			result, err := statement.Exec(messageBody.Message, messageBody.Content)
			failOnError(err, "Failed to execute query")

			// Get and print Inserted ID
			id, err := result.LastInsertId()
			failOnError(err, "Failed to retrieve lastInsertId")
			fmt.Println(id)
		}
	}()

	// Log output of channel
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}