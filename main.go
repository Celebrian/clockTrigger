package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var latest int64

func main() {
	go count10()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server Started")
}

func count10() {
	for {
		go checkNew()
		time.Sleep(time.Minute)
	}
}

func checkNew() {
	resp, err := http.Get("http://localhost:8080/paragliding/api/ticker/latest")
	if err != nil {
		fmt.Println("Failed to get latest ticker")
		return
	}
	defer resp.Body.Close()

	latestNow, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read response Body")
		return
	}
	latestNowString, err := strconv.Atoi(string(latestNow))
	if err != nil {
		fmt.Println("Could not convert response body to int")
		return
	}
	latest64 := int64(latestNowString)

	if latest64 != latest {
		type payloadStruct struct {
			PayloadString string `json:"content"`
		}

		payload := payloadStruct{
			PayloadString: "New tracks added to the database",
		}

		jsonPayload, _ := json.Marshal(payload)
		_, err := http.Post(os.Getenv("WEBHOOK"), "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("Could not post to webhook")
			return
		}
		latest = latest64
	}
}
