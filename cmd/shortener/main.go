package main

import (
	"errors"
	"log"
	"net/http"
	"new-shortner/internal/app"
	"new-shortner/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func a() {
	//TODO
}
