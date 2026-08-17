package processors

import (
	"goka_chain/common"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/lovoo/goka"
)

// purchasesEmitter — эмиттер, который генерирует данные в топик purchases
func PurchasesEmitter(brokers []string, topicOrders goka.Stream) {
	e, err := goka.NewEmitter(brokers, topicOrders, new(common.JsonCodec[Order]))
	if err != nil {
		log.Fatal(err)
	}

	defer e.Finish()

	for {
		time.Sleep(1 * time.Second)

		up := Order{
			UserID:      rand.Int64N(10),            // Случайный идентификатор пользователя в диапазоне [0, 10)
			OrderID:     rand.Int64(),               // Случайный идентификатор покупки
			OrderAmount: 1000 + rand.Int64N(90_000), // Случайный сумма покупки в диапазоне [1_000, 100_000)
		}

		if err = e.EmitSync(strconv.FormatInt(up.UserID, 10), up); err != nil {
			log.Fatal(err)
		}
		log.Printf("Новая покупка пользователя %d на сумму %d\n", up.UserID, up.OrderAmount)
	}
}
