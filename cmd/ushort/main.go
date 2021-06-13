package main

import (
	"fmt"
	"log"

	"github.com/AleksandrMac/ushort/pkg/config"
)

func main() {
	var (
		conf *config.Config
		err  error
	)

	conf, err = config.New("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
}
