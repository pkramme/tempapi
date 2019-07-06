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

	AddNewRoomString     string = "INSERT INTO rooms (name) VALUES ($1) RETURNING id"
	AddNewTempdataString string = "INSERT INTO tempdata (roomid, time, temperature) VALUES ($1, $2, $3)"
	GetRoomByNameString  string = "SELECT id FROM rooms WHERE name = $1"
)

var (
	AccessToken string

	TimezoneLocation *time.Location

	AddNewRoomSQL     *sql.Stmt
	AddNewTempdataSQL *sql.Stmt
	GetRoomByNameSQL  *sql.Stmt
)

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
		unixutctimeval := time.Unix(timeint, 0)
		timeval := unixutctimeval.In(TimezoneLocation)

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

		// TODO: Implement a map where the room and the last received temperature is saved. Compare the current temperature vs the old temperature for every request. Write only to the database when the temperature has changed. IO is expensive.

		var roomid int64
		err = GetRoomByNameSQL.QueryRow(room).Scan(&roomid)
		if err != nil {
			fmt.Println(err)
		}
		if roomid == 0 {
			fmt.Println(room, "doesnt exist, creating it.")
			err := AddNewRoomSQL.QueryRow(room).Scan(&roomid)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		_, err = AddNewTempdataSQL.Exec(roomid, timeval, temp)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
	Timezone := os.Getenv("timezone")

	var err error
	TimezoneLocation, err = time.LoadLocation(Timezone)
	if err != nil {
		panic(err)
	}

	fmt.Println("connecting to database")

	// Setup database connection
	DatabaseConnectionString := fmt.Sprintf("dbname=%s user=%s host=%s password=%s", DatabaseName, DatabaseUser, DatabaseAddress, DatabasePassword)
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
	_, err = db.Exec(SetupRoomTable)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(SetupTempdataTable)
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

	AddNewTempdataSQL, err = db.Prepare(AddNewTempdataString)
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
