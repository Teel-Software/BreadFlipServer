package main

import (
	"hleb_flip/internal/router"
)

func main() {
	r := router.NewRouter("0.0.0.0", "8080")
	r.StartRouter()
}
