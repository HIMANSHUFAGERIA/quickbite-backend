package main

import (
	"fmt"
	"log"
	"net/http"

	"quickbite/config"
	"quickbite/db"
	"quickbite/internal/handler"
	"quickbite/internal/middleware"
)

func main() {
	cfg := config.Load()

	db.Connect(cfg)
	defer db.DB.Close()

	mux := handler.NewRouter(cfg)

	wrappedMux := middleware.CORS(cfg)(middleware.Logger(mux))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 QuickBite API running on http://localhost%s", addr)
	log.Printf("📦 Environment: %s", cfg.Environment)

	if err := http.ListenAndServe(addr, wrappedMux); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
