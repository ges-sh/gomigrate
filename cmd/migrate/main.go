package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var name string
var dir string

func init() {
	flag.StringVar(&name, "name", "", "Name of the migration")
	flag.StringVar(&dir, "dir", "", "Directory for the migration")
	flag.Parse()
}

func main() {
	os.Exit(run())
}

func run() int {
	if name == "" {
		log.Println("migrate: name cannot be empty")
		return 1
	}

	timestamp := time.Now().Unix()

	mFormat := "%d_%s.%s.sql"

	_, err := os.Create(fmt.Sprintf(mFormat, timestamp, name, "up"))
	if err != nil {
		log.Printf("migrate: can't create up migration: %v\n", err)
		return 1
	}
	_, err = os.Create(fmt.Sprintf(mFormat, timestamp, name, "down"))
	if err != nil {
		log.Printf("migrate: can't create down migration: %v\n", err)
		return 1
	}

	return 0
}
