package main

import (
	"fmt"
	"net/http"

	"github.com/jonggulee/go-login.git/src/api/controller"
)

func main() {
	fmt.Printf("Hello, Ggomo\n\n")
	http.ListenAndServe(":8080", controller.NewRouter())
}
