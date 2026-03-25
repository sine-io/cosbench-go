package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sine-io/cosbench-go/internal/app"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "http listen address")
	dataDir := flag.String("data-dir", "data", "snapshot data directory")
	viewDir := flag.String("view-dir", "web/templates", "html template directory")
	flag.Parse()

	application, err := app.New(app.Config{DataDir: *dataDir, ViewDir: *viewDir})
	if err != nil {
		log.Fatalf("bootstrap server: %v", err)
	}

	fmt.Printf("cosbench-go server listening on %s\n", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, application.Handler); err != nil {
		log.Fatal(err)
	}
}
