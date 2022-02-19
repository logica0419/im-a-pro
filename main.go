package main

import (
	"log"

	"github.com/logica0419/im-a-pro/router"
)

func main() {
	r, err := router.NewRouter()
	if err != nil {
		log.Panic(err)
	}

	r.Run()
}
