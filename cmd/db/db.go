package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// DBUser is the database username
var DBUser = os.Getenv("DB_USER")

// DBName is the name of the database
var DBName = os.Getenv("DB_NAME")

// DBPassword is the database password
var DBPassword = os.Getenv("DB_PASSWORD")

// Message represents the structure of the message to insert
type Message struct {
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"time"`
	DeviceID    string    `json:"sensor_id"`
}

func Connect() error {
	// connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", DBUser, DBName, DBPassword)
	connStr := fmt.Sprintf("postgres://%s:%s@192.168.178.55/%s?sslmode=disable", DBUser, DBPassword, DBName)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	return nil
}

// InsertMessage inserts a Message into the database
func InsertMessage(msg Message) {
	// insert data into the database
	_, err := DB.Exec("INSERT INTO sensor_readings (temperature, humidity, reading_time, sensor_id) VALUES ($1, $2, $3, $4)", msg.Temperature, msg.Humidity, msg.Timestamp, msg.DeviceID)
	if err != nil {
		log.Fatalf("Failed to insert message: %s", err)
	}
	log.Println("Message inserted successfully")
}
