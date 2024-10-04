package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func notifyUser(bid Bid) {
	body, _ := json.Marshal(bid)
	req, _ := http.NewRequest(http.MethodPut, bid.Update, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	failOnError(err, "Error sending the HTTP request")
}
