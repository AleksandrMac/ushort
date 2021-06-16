package main

import (
	"log"

	"github.com/AleksandrMac/ushort/pkg/config"
)

var (
	c   *config.Config
	err error
)

func main() {
	c, err = config.New("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
}
