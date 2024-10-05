package main

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	h "auctions/helper"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	h.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	go h.RecieveMessages(conn, "create_auctions", handleAuctionCreate)
	go h.RecieveMessages(conn, "search_auctions", handleSearchAuctions)
	go h.RecieveDelayMessages(conn, "close_auctions", handleCloseAuction)

	select {}
}

func handleAuctionCreate(d amqp.Delivery, channel *amqp.Channel) {
	var auction h.Auction
	json.Unmarshal(d.Body, &auction)

	auction = insertAuction(auction)
	acutionData, _ := json.Marshal(auction)

	h.ReplyMessage(channel, d, acutionData)
}

func handleSearchAuctions(d amqp.Delivery, channel *amqp.Channel) {
	var params h.AuctionSearchParams
	json.Unmarshal(d.Body, &params)

	result := searchAuctions(params)
	resultData, _ := json.Marshal(result)

	h.ReplyMessage(channel, d, resultData)
}

func handleCloseAuction(d amqp.Delivery, channel *amqp.Channel) {
	var auction h.Auction
	json.Unmarshal(d.Body, &auction)

	bids := getAuctionBids(auction)

	for _, bid := range bids {
		if bid.Status == "best" {
			h.UpdateStatus(Connect(), bid, "won")
		} else {
			h.UpdateStatus(Connect(), bid, "lost")
		}
	}
}
