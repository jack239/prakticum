package processors

import (
	"context"
	"goka_chain/common"
	"log"

	"github.com/lovoo/goka"
)

type (
	Order        = common.Order
	UserSum      = common.UserSum
	UserSumCodec = common.JsonCodec[UserSum]
	OrderCodec   = common.JsonCodec[Order]
)

var (
	groupUsersOrderSum goka.Group = "users-sum-group"
)

// sumProcessor считает сумму всех заказов
func SumProcessor(brokers []string, inputTopic, outputTopic goka.Stream) {
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
		ctx.Emit(outputTopic, ctx.Key(), userSum)
	}

	g := goka.DefineGroup(groupUsersOrderSum,
		goka.Input(inputTopic, new(OrderCodec), processFunc),
		goka.Persist(new(UserSumCodec)),
		goka.Output(outputTopic, new(UserSumCodec)),
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
