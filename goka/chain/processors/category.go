package processors

import (
	"context"
	"goka_chain/common"
	"log"

	"github.com/lovoo/goka"
)

var (
	groupUsersCategory goka.Group = "users-category-group"
)

// categoryProcessor присваивает категорию пользователю
func CategoryProcessor(brokers []string, input, output goka.Stream) {
	processFunc := func(ctx goka.Context, msg interface{}) {
		var (
			userSum common.UserSum
			ok      bool
		)
		if userSum, ok = msg.(common.UserSum); !ok {
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

		ctx.Emit(output, ctx.Key(), common.UserCategory{Category: category})
	}

	g := goka.DefineGroup(groupUsersCategory,
		goka.Input(input, new(common.JsonCodec[common.UserSum]), processFunc),
		goka.Output(output, new(common.JsonCodec[common.UserCategory])),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	if err = p.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
