package common

import (
	"encoding/json"
	"fmt"
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
