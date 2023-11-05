package main

import (
	"os"

	"github.com/jonggulee/como-backend/src/app"
)

func main() {
	app.App.Run(os.Args)
}
