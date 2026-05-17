package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zrgst/arx-dash/internal/arx"
	"github.com/zrgst/arx-dash/internal/config"
	"github.com/zrgst/arx-dash/internal/web"
)

func main() {
	cfg := config.Load()

	if cfg.ARXBaseURL == "" {
		log.Fatal("Missing ARX_BASE_URL")
	}

	if cfg.ARXUsername == "" {
		log.Fatal("Missing ARX_USERNAME")
	}

	arxClient := arx.NewClient(
		cfg.ARXBaseURL,
		cfg.ARXUsername,
		cfg.ARXPassword,
		cfg.ARXAllowSelfSigned,
	)

	server := web.NewServer(arxClient)

	fmt.Printf("ARX Dashboard Go listening on http://localhost%s\n", cfg.Addr)

	if err := http.ListenAndServe(cfg.Addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
