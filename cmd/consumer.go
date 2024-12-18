/*
Copyright © 2024 Skrypnyk Yuriy <skrypnyk81@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Skrypnyk81/manager_esp_sensor/cmd/db"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var (
	// rabbitmq credentials
	RabbitMQUser     = os.Getenv("RABBITMQ_USER")
	RabbitMQPassword = os.Getenv("RABBITMQ_PASSWORD")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:     "consumer",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Connessione al database
		err := db.Connect()
		failOnError(err, "Failed to connect to the database")

		// Connessione al server RabbitMQ
		conn, err := amqp.Dial("amqp://" + RabbitMQUser + ":" + RabbitMQPassword + "@192.168.178.55:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		// Creazione del canale
		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		// Dichiarazione della coda
		queueName := "esp8266_amqp"
		q, err := ch.QueueDeclare(
			queueName, // Nome della coda
			true,      // Durable
			false,     // Delete when unused
			false,     // Exclusive
			false,     // No-wait
			nil,       // Arguments
		)
		failOnError(err, "Failed to declare a queue")

		// Associazione dell'exchange `amq.topic` alla coda `esp8266_amqp`
		err = ch.QueueBind(
			q.Name,      // Nome della coda
			"#",         // Routing key
			"amq.topic", // Nome dell'exchange
			false,       // No-wait
			nil,         // Arguments
		)
		failOnError(err, "Failed to bind the queue")

		// Consumo dei messaggi dalla coda
		msgs, err := ch.Consume(
			q.Name, // Nome della coda
			"",     // Consumer
			true,   // Auto-ack
			false,  // Exclusive
			false,  // No-local
			false,  // No-wait
			nil,    // Args
		)
		failOnError(err, "Failed to register a consumer")

		// Canale per segnalare la fine del programma
		forever := make(chan bool)

		// Goroutine per ricevere i messaggi
		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				// unmarshal the message
				var message db.Message
				err := json.Unmarshal(d.Body, &message)
				failOnError(err, "Failed to unmarshal the message")
				// insert the message into the database
				db.InsertMessage(message)
			}
		}()

		log.Printf("Waiting for messages. To exit press CTRL+C")
		<-forever
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
