package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	SetupRoomTable     string = "CREATE TABLE IF NOT EXISTS rooms (id SERIAL PRIMARY KEY NOT NULL, name varchar(255) NOT NULL)"
	SetupTempdataTable string = "CREATE TABLE IF NOT EXISTS tempdata (id SERIAL PRIMARY KEY NOT NULL, roomid SERIAL REFERENCES rooms(id) NOT NULL, time TIMESTAMP NOT NULL, temperature FLOAT(32) NOT NULL)"

	AddNewRoomString     string = "INSERT INTO rooms (name) VALUES ($1)"
	AddNewTempdataString string = "INSERT INTO tempdata (roomid, time, temperature) VALUES ($1, $2, $3)"
	GetRoomByNameString  string = "SELECT id FROM rooms WHERE name = $1"
)

var AccessToken string

func TempPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		q := r.URL.Query()

		// Auth
		token := q.Get("token")
		if token == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		if token != AccessToken {
			http.Error(w, "Wrong token", http.StatusForbidden)
			return
		}

		// Get time
		timestring := q.Get("time")
		if timestring == "" {
			http.Error(w, "No time given", http.StatusBadRequest)
			return
		}
		timeint, err := strconv.ParseInt(timestring, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		timeval := time.Unix(timeint, 0)
		fmt.Println(timeval)

		// Get room
		room := q.Get("room")
		if room == "" {
			http.Error(w, "No room given", http.StatusBadRequest)
			return
		}

		// Get temp
		tempstring := q.Get("temp")
		if tempstring == "" {
			http.Error(w, "No temp given", http.StatusBadRequest)
			return
		}
		temp, err := strconv.ParseFloat(tempstring, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Room:", room)
		fmt.Println("Time:", timeval)
		fmt.Println("Temperature:", float32(temp))

		// TODO: Implement a map where the room and the last received temperature is saved. Compare the current temperature vs the old temperature for every request. Write only to the database when the temperature has changed. IO is expensive.

		// TODO: Implement writing to database
		//ok, lets go.
		// first, we need to check if the room exists and get its id.
		//    if it doesnt, create it and get its id.
		// then, we insert the data with the foreign key set to the id, the time provided and the actual temperature

	} else {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	TLSKey := os.Getenv("tlskey")
	TLSCert := os.Getenv("tlscert")
	AccessToken = os.Getenv("token")
	DatabaseAddress := os.Getenv("databaseaddress")
	DatabasePassword := os.Getenv("databasepassword")
	DatabaseUser := os.Getenv("databaseuser")
	DatabaseName := os.Getenv("databasename")

	fmt.Println("connecting to database")

	// Setup database connection
	DatabaseConnectionString := fmt.Sprintf("dbname=%s user=%s host=%s password=%s", DatabaseUser, DatabaseAddress, DatabasePassword)
	db, err := sql.Open("postgres", DatabaseConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection to database successfull")

	fmt.Println("beginning database setup")
	_, err := db.Exec(SetupRoomTable)
	if err != nil {
		panic(err)
	}
	_, err := db.Exec(SetupTempdataTable)
	if err != nil {
		panic(err)
	}
	fmt.Println("database setup done")
	fmt.Println("preparing sql statements")

	AddNewRoomSQL, err = db.Prepare(AddNewRoomString)
	if err != nil {
		panic(err)
	}
	defer AddNewRoomSQL.Close()

	AddNewTempdataSQL, err = db.Prepare(AddNewTempdataSQL)
	if err != nil {
		panic(err)
	}
	defer AddNewTempdataSQL.Close()

	GetRoomByNameSQL, err = db.Prepare(GetRoomByNameString)
	if err != nil {
		panic(err)
	}
	defer GetRoomByNameSQL.Close()

	fmt.Println("sql statement preparation done")

	fmt.Println("starting webserver")
	http.HandleFunc("/", TempPost)
	log.Fatal(http.ListenAndServeTLS(":8443", TLSCert, TLSKey, nil))
}
