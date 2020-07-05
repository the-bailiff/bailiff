package main

import (
	"net/http"

	"github.com/the-bailiff/bailiff/internal/app"
)

func main() {
	a, err := app.InitApp()
	if err != nil {
		panic(err)
	}

	http.Handle("/", a.Handler)

	if err := http.ListenAndServe(":"+a.Port, nil); err != nil {
		panic(err)
	}
}
