package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 7080
	defaultHost    = "0.0.0.0"
	appName        = "sub2api"
	appVersion     = "1.0.0"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", getEnvInt("PORT", defaultPort), "Port to listen on")
	host := flag.String("host", getEnv("HOST", defaultHost), "Host to bind to")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handler.HealthCheck)

	// Subscription conversion endpoints
	mux.HandleFunc("/sub", handler.ConvertSubscription)
	mux.HandleFunc("/clash", handler.ConvertToClash)
	mux.HandleFunc("/singbox", handler.ConvertToSingBox)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Starting %s v%s on %s", appName, appVersion, addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// getEnvInt returns the integer value of an environment variable or a default value.
func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
