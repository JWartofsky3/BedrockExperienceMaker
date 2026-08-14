// Command server exposes the small HTTP API consumed by the Vite application.
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JWartofsky3/BedrockExperienceMaker/api/cmd/server/addons/operations"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN must be set; see .env.example")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()
	log.Print("Connecting to MySQL...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	log.Print("Connected to MySQL")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /v1/addons", operations.ListAddons(db))
	mux.HandleFunc("GET /v1/addons/{name}", operations.GetAddon(db))
	mux.HandleFunc("DELETE /v1/addons/{name}", operations.DeleteAddon(db))
	address := os.Getenv("HTTP_ADDRESS")
	if address == "" {
		address = "127.0.0.1:8080"
	}
	log.Printf("Add-on API listening on http://%s", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
