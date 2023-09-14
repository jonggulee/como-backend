package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/jonggulee/go-login.git/src/api/controller"
	"github.com/jonggulee/go-login.git/src/logger"
)

const (
	port = 8080
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Logger.Fatal("Error loading .env file")
	}
}

func main() {
	startupMessage()
	http.ListenAndServe(fmt.Sprintf(":%d", port), controller.NewRouter())
}

func startupMessage() {
	fmt.Println("==================================================")
	fmt.Println("                GGOMO API SERVER")
	fmt.Println("==================================================")
	fmt.Printf("Server is starting...\n\n")
	fmt.Printf("Version:		%s\n", "v0.1")
	fmt.Printf("Listening port: 	%d\n", port)
	fmt.Printf("Started at: 		%s\n", time.Now().Format(time.RFC3339))
	fmt.Println("==================================================")
}
