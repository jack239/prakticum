package main

import (
	"goka_chain/processors"

	"github.com/lovoo/goka"
)

var (
	brokers = []string{"localhost:9092"}

	topicOrders        goka.Stream = "orders"
	topicUsersSum      goka.Stream = "users.sum"
	topicUsersCategory goka.Stream = "users.category"
)

func main() {
	go processors.PurchasesEmitter(brokers, topicOrders)
	go processors.SumProcessor(brokers, topicOrders, topicUsersSum)
	go processors.CategoryProcessor(brokers, topicUsersSum, topicUsersCategory)
	go processors.LoggerProcessor(brokers, topicUsersCategory)

	select {} // Блокируем main, чтобы горутины работали
}
