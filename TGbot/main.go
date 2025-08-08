package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)
var TOKEN string

func init() {
	godotenv.Load()
	TOKEN = os.Getenv("TOKEN")
}

func main() {
	fmt.Println(TOKEN)
}
