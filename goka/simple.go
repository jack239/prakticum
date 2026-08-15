package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type User struct {
	Name string `json:"name"`
}

type UserCodec struct{}

func (jc *UserCodec) Encode(value interface{}) ([]byte, error) {
	// Ваш код тут
	if user, ok := value.(User); ok {
		return []byte(user.Name), nil
	}
	return nil, fmt.Errorf("Illiegal type")
}

func (jc *UserCodec) Decode(data []byte) (interface{}, error) {
	// Ваш код тут
	return nil, nil
}

var (
	brokers             = []string{"localhost:9092"} // Адрес брокера
	input   goka.Stream = "input"                    // Топик с исходными данными
	output  goka.Stream = "output"                   // Топик с результатом
	group   goka.Group  = "group"                    // Имя группы
)

// Генерирует сообщения (строки) и отправляет их в топик input каждую секунду
func runEmitter() {
	fmt.Println("Emitor started")
	emitter, err := goka.NewEmitter(brokers, input, new(codec.String))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitter.Finish()

	var counter int
	for {
		time.Sleep(1 * time.Second)
		err = emitter.EmitSync("key", fmt.Sprintf("Value #%d", counter))
		if err != nil {
			log.Fatalf("error emitting message: %v", err)
		}
		log.Printf("[emitter] Сообщение #%d отправлено\n", counter)
		counter++
	}
}

func runProcss() {
	fmt.Println("Processor started")
	upperCaseFunc := func(ctx goka.Context, msg interface{}) {
		// Обработчик — преобразует значения сообщений в upper-case и отправляет их в output
		log.Printf("[processor] Получено сообщение: key = %s, value = %s", ctx.Key(), msg)

		if strMsg, ok := msg.(string); ok {
			upper := strings.ToUpper(strMsg)

			// Отправляем сообщения в output
			ctx.Emit(output, ctx.Key(), upper)
			log.Printf("[processor] Сообщение обработано: key = %s, new_value = %s", ctx.Key(), upper)
		}
	}

	// Создание группы
	g := goka.DefineGroup(group,
		goka.Input(input, new(codec.String), upperCaseFunc),
		goka.Output(output, new(codec.String)),
	)
	// Создание группы

	// Создание процессора
	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	// Запуск процессора в отдельной горутине
	ctx := context.Background()
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = p.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()
}

func main() {
	fmt.Println("Main started")
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Emitor started")
		runEmitter()
		fmt.Println("Emitor finished")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Processor started")
		runProcss()
		fmt.Println("Processor finished")
	}()
	wg.Wait()
	fmt.Println("Main finished")
}
