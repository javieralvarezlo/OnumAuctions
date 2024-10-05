package main

import (
	"encoding/json"
	"fmt"
	"time"

	h "auctions/helper"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	h.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	go h.RecieveMessages(conn, "create_bids", handleBidCreate)
	go h.RecieveMessages(conn, "search_bids", handleBidsSearch)

	select {}
}

func handleBidCreate(d amqp.Delivery, channel *amqp.Channel) {
	var bid h.Bid
	json.Unmarshal(d.Body, &bid)

	auction, err := getBidAuction(bid)

	if err != nil {
		bid.BidID = "-1"
		bid.Status = "The selected auction does not exist"

		bidData, _ := json.Marshal(bid)
		h.ReplyMessage(channel, d, bidData)
		return
	}

	if time.Now().UnixMilli() < auction.BidStartTime {
		bid.BidID = "-1"
		bid.Status = "The auction has not started yet"

		bidData, _ := json.Marshal(bid)
		h.ReplyMessage(channel, d, bidData)
		return
	}

	if time.Now().UnixMilli() > auction.BidEndTime {
		bid.BidID = "-1"
		bid.Status = "The auction has already finished"

		bidData, _ := json.Marshal(bid)
		h.ReplyMessage(channel, d, bidData)
		return
	}

	if bid.Value < int(auction.StartValue) {
		bid.BidID = "-1"
		bid.Status = fmt.Sprint("The initial value of the auction is ", auction.StartValue)

		bidData, _ := json.Marshal(bid)
		h.ReplyMessage(channel, d, bidData)
		return
	}

	bid.Status = "processing"

	bid = insertBid(bid)
	bidData, _ := json.Marshal(bid)

	defer processBid(bid)

	h.ReplyMessage(channel, d, bidData)
}

func handleBidsSearch(d amqp.Delivery, channel *amqp.Channel) {
	var params h.BidSearchParams
	json.Unmarshal(d.Body, &params)

	result := searchBids(params)
	resultData, _ := json.Marshal(result)

	h.ReplyMessage(channel, d, resultData)
}
