package main

import (
	"fmt"
	"net/http"
)

// сюда код писать не надо
// TODO: Регистрация не должна возвращать токен. Аутентификация происходит после логина
// TODO: Все ошибки должны быть в форме констант

func main() {
	addr := "127.0.0.1:8080" // TODO: Вернуть листинг всех интерфейсов :8080
	h := GetApp()
	fmt.Println("start server at", addr)
	http.ListenAndServe(addr, h)
}
