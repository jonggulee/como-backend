package main

import (
	"os"

	"github.com/jonggulee/go-login.git/src/app"
)

func main() {
	app.App.Run(os.Args)
}
