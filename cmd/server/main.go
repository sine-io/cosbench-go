package main

import (
	"context"
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
	mode := flag.String("mode", string(app.ModeCombined), "runtime mode: controller-only, driver-only, combined")
	driverSharedToken := flag.String("driver-shared-token", "", "shared bearer token for /api/driver write endpoints; falls back to COSBENCH_DRIVER_SHARED_TOKEN")
	controllerURL := flag.String("controller-url", "", "controller base URL for driver-only mode; falls back to COSBENCH_CONTROLLER_URL")
	driverName := flag.String("driver-name", "", "driver name for driver-only mode; falls back to COSBENCH_DRIVER_NAME")
	flag.Parse()

	application, err := app.New(app.Config{
		DataDir: *dataDir,
		ViewDir: *viewDir,
		Mode: app.Mode(*mode),
		DriverSharedToken: *driverSharedToken,
		ControllerURL: *controllerURL,
		DriverName: *driverName,
	})
	if err != nil {
		log.Fatalf("bootstrap server: %v", err)
	}
	if err := application.StartBackground(context.Background()); err != nil {
		log.Fatalf("background runtime: %v", err)
	}

	fmt.Printf("cosbench-go server listening on %s\n", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, application.Handler); err != nil {
		log.Fatal(err)
	}
}
