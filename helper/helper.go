package helper

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NotifyUser(bid Bid) {
	body, _ := json.Marshal(bid)
	req, _ := http.NewRequest(http.MethodPut, bid.Update, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	FailOnError(err, "Error sending the HTTP request")
}
