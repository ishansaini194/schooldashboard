package main

import (
	"log"

	"github.com/ishansaini194/dashboard/internal/app"
)

func main() {
	srv := app.New()
	log.Fatal(srv.Start(":8080"))
}
