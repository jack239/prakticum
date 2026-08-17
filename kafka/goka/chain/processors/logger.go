package processors

import (
	"context"
	"goka_chain/common"
	"log"

	"github.com/lovoo/goka"
)

type (
	UserCategory = common.UserCategory
)

var (
	groupLogger goka.Group = "logger-group"
)

// loggerProcessor просто логирует финальную категорию пользователя
func LoggerProcessor(brokers []string, topicUsersCategory goka.Stream) {
	processFunc := func(ctx goka.Context, msg interface{}) {
		if userCategory, ok := msg.(UserCategory); ok {
			log.Printf("Категория пользователя %s = %s\n", ctx.Key(), userCategory.Category)
		}
	}

	g := goka.DefineGroup(groupLogger,
		goka.Input(topicUsersCategory, new(common.JsonCodec[UserCategory]), processFunc),
	)

	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatal(err)
	}
	if err = p.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
