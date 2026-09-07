package main

import (
	"fmt"
	"os"
	"text/template"
)

// Person структура, данные которой будут подставляться в шаблон
type Person struct {
	Name string
	Age  int
}

func alice() {
	alicePerson := Person{"Алиса", 7}
	aliceTemplate := "Это {{ .Name }} и ей {{ .Age }} лет."

	// создаем новый шаблон и парсим его содержимое, тем самым подготавливая его к дальнейшему использованию
	readyTemplate, err := template.New("test").Parse(aliceTemplate)
	if err != nil {
		fmt.Println("Ошибка при создании шаблона")
		return
	}

	// выполняем шаблон и выводим его стандартный поток
	err = readyTemplate.Execute(os.Stdout, alicePerson)
	if err != nil {
		fmt.Println("Ошибка при выполнении шаблона")
		return
	}
}

func main() {
	weekDays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	weekTemplate := "Все дни недели: {{ range $index, $element := .}}{{ if $index }}, {{ end }}{{$element}}{{end}}"
	readyWeekTemplate, err := template.New("test").Parse(weekTemplate)

	if err != nil {
		panic(err)
	}

	readyWeekTemplate.Execute(os.Stdout, weekDays)
}
