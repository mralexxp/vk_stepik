package main

import (
	"testing"
)

func TestNewBroadcast(t *testing.T) {
	inputChan := make(chan *Event)
	bc := NewBroadcast(inputChan)

	// Проверяем запуск службы
	firstChan, secondChan, deleteChan := bc.Subscribe(), bc.Subscribe(), bc.Subscribe()

	// Проверяем количество подписоты
	if len(bc.subscribers) != 3 {
		t.Errorf("len(bc.subscribers)=%d, want 3", len(bc.subscribers))
	}

	bc.Unsubscribe(deleteChan)

	if len(bc.subscribers) != 2 {
		t.Errorf("len(bc.subscribers)=%d, want 2", len(bc.subscribers))
	}

	bc.Unsubscribe(firstChan)
	bc.Unsubscribe(secondChan)

	//testEvents := []Event{
	//	{
	//		Timestamp:     1,
	//		Consumer:      "123",
	//		Method:        "TestMethod",
	//		Host:          "127.0.1.2:12345",
	//	},
	//	{
	//		Timestamp:     2,
	//		Consumer:      "321",
	//		Method:        "TestMethodTwo",
	//		Host:          "127.0.1.2:12345",
	//	},
	//	{
	//		Timestamp:     3,
	//		Consumer:      "321123",
	//		Method:        "TestMethodThree",
	//		Host:          "127.0.1.2:12345",
	//	},
	//}

	// TODO: Проверяем удаление подписчика

	// TODO: Проверяем вентилятор тремя каналами
}
