package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
var TOKEN string

func init() {
	godotenv.Load()
	TOKEN = os.Getenv("TOKEN")
}

func main() {
	if TOKEN == "" {
		log.Fatal("")
	}


}
