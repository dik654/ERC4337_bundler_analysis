package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	url := "http://localhost:5000/publish"
	topics := []string{"foobarbaz", "topic_2", "topic_3"}

	for i := 0; i < 100; i++ {
		topic := topics[rand.Intn(len(topics))]
		payload := []byte(fmt.Sprintf("foobarbaz_%d", i))

		resp, err := http.Post(url+"/"+topic, "application/json", bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal("status code is not 200")
		}
	}
}
