package main

import (
	"log"

	"goFinal/pkg/server"
)

func main() {

	if err := server.Run(); err != nil {
		log.Println(err)
	}
}
