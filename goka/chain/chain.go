package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/lovoo/goka"
)

// Order описывает сообщение в топике с заказами
type Order struct {
	UserID      int64 `json:"user_id"`
	OrderID     int64 `json:"order_id"`
	OrderAmount int64 `json:"order_amount"`
}

// UserSum сумма всех заказов пользователя
type UserSum struct {
	Total int64 `json:"total"`
}

// UserCategory категория пользователя
type UserCategory struct {
	Category string `json:"category"`
}

// JsonCodec общий кодек, который предназначен для сериализации и десериализации в json
type JsonCodec[T any] struct{}

func (jc JsonCodec[T]) Encode(value interface{}) ([]byte, error) {
	if user, ok := value.(T); ok {
		return json.Marshal(user)
	}

	return nil, fmt.Errorf("illegal type: %T", value)
}

func (jc JsonCodec[T]) Decode(data []byte) (interface{}, error) {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return t, nil
}

var (
	brokers = []string{"localhost:9092"}

	topicOrders        goka.Stream = "orders"
	topicUsersSum      goka.Stream = "users.sum"
	topicUsersCategory goka.Stream = "users.category"

	groupUsersOrderSum goka.Group = "users-sum-group"
	groupUsersCategory goka.Group = "users-category-group"
	groupLogger        goka.Group = "logger-group"
)

// purchasesEmitter — эмиттер, который генерирует данные в топик purchases
func purchasesEmitter() {
	e, err := goka.NewEmitter(brokers, topicOrders, new(JsonCodec[Order]))
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

// sumProcessor считает сумму всех заказов
func sumProcessor() {
	processFunc := func(ctx goka.Context, msg interface{}) {
		var (
			order Order
			ok    bool
		)
		if order, ok = msg.(Order); !ok {
			log.Printf("illegal type: %T", msg)
			return // Не останавливаем процессор
		}

		// Считываем текущее значение для пользователя
		var userSum UserSum
		currentSum := ctx.Value() // Значение для ключа сообщения — а оно совпадает с идентификатором пользователя
		if currentSum != nil {
			userSum = currentSum.(UserSum)
		}

		// И добавляем сумму к этому пользователю
		userSum.Total += order.OrderAmount
		ctx.SetValue(userSum)
		log.Printf("Текущая сумма заказов пользователя %s: %d\n", ctx.Key(), userSum.Total)

		// Отправляем обновленную сумму в следующий топик
		ctx.Emit(topicUsersSum, ctx.Key(), userSum)
	}

	g := goka.DefineGroup(groupUsersOrderSum,
		goka.Input(topicOrders, new(JsonCodec[Order]), processFunc),
		goka.Persist(new(JsonCodec[UserSum])),
		goka.Output(topicUsersSum, new(JsonCodec[UserSum])),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Stop()

	if err = p.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// categoryProcessor присваивает категорию пользователю
func categoryProcessor() {
	processFunc := func(ctx goka.Context, msg interface{}) {
		var (
			userSum UserSum
			ok      bool
		)
		if userSum, ok = msg.(UserSum); !ok {
			log.Printf("illegal type: %T", msg)
			return // Не останавливаем процессор
		}

		var category string
		switch {
		case userSum.Total >= 1_000_000: // Если сумма покупок больше 1_000_000, то категория gold
			category = "gold"
		case userSum.Total >= 500_000: // Если сумма заказов больше 500_000, то категория silver
			category = "silver"
		default: // Иначе — категория bronze
			category = "bronze"
		}

		ctx.Emit(topicUsersCategory, ctx.Key(), UserCategory{Category: category})
	}

	g := goka.DefineGroup(groupUsersCategory,
		goka.Input(topicUsersSum, new(JsonCodec[UserSum]), processFunc),
		goka.Output(topicUsersCategory, new(JsonCodec[UserCategory])),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	if err = p.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// loggerProcessor просто логирует финальную категорию пользователя
func loggerProcessor() {
	processFunc := func(ctx goka.Context, msg interface{}) {
		if userCategory, ok := msg.(UserCategory); ok {
			log.Printf("Категория пользователя %s = %s\n", userCategory.Category, userCategory.Category)
		}
	}

	g := goka.DefineGroup(groupLogger,
		goka.Input(topicUsersCategory, new(JsonCodec[UserCategory]), processFunc),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	if err = p.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func main() {
	go purchasesEmitter()
	go sumProcessor()
	go categoryProcessor()
	go loggerProcessor()

	select {} // Блокируем main, чтобы горутины работали
}
