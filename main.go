package main

import (
	"os"

	"github.com/jonggulee/go-login.git/src/app"
)

// func init() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		logger.Logger.Fatal("Error loading .env file")
// 	}

// 	fmt.Println(godotenv.Read(".env"))
// }

func main() {
	app.App.Run(os.Args)

}
