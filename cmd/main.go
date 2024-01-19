package main

import (
	"hleb_flip/internal/router"
	"log"
	"os"
)

func main() {
	os.Create("/tmp/breadlog.txt")
	f, err := os.Open("/tmp/breadlog.txt")
	if err != nil {
		log.Default().Println("failed to open log file")
	} else {
		log.Default().SetOutput(f)

	}
	r := router.NewRouter("0.0.0.0", "8080")
	r.StartRouter()
}
