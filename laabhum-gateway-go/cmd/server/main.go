package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/Mukilan-T/laabhum-gateway-go/routes"
)

// Order struct represents an order in the system.
type Order struct {
	ID                string               `json:"id"`
	Symbol            string               `json:"symbol"`
	Quantity          int                  `json:"quantity"`
	Price             float64              `json:"price"`
	Side              string               `json:"side"` 
	Type              string               `json:"type"`
	Status            string               `json:"status"` 
	StopPrice         float64              `json:"stop_price,omitempty"` 
	Strategy          string               `json:"strategy"`
	RiskPercentage    float64              `json:"risk_percentage"` 
	StopLossActivated bool                 `json:"stop_loss_activated"` 
	TakeProfit        float64              `json:"take_profit"` 
	CreatedAt         int64                `json:"created_at"` 
	ExpiresAt         time.Time            `json:"expires_at,omitempty"` 
	ParentID          string               `json:"parent_id"`
}

func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("Failed to load configuration")
	}
	fmt.Printf("Loaded OMS address: %s\n", cfg.Oms.BaseURL)

	logLevel := cfg.LogLevel
	customLogger := logger.New(logLevel)

	stdLogger := log.New(customLogger.Writer(), "", log.LstdFlags)

	omsClient := oms.NewClient(cfg.Oms.BaseURL)
	if omsClient == nil {
		stdLogger.Fatalf("Failed to create OMS client")
	}

	// Fetch existing orders
	ordersData, err := omsClient.GetOrders()
	if err != nil {
		stdLogger.Fatalf("Failed to get orders: %v", err)
	}

	var orders []Order
	if err := json.Unmarshal(ordersData, &orders); err != nil {
		stdLogger.Fatalf("Failed to unmarshal orders: %v", err)
	}

	for _, order := range orders {
		stdLogger.Printf("Order ID: %s, Status: %s", order.ID, order.Status)
	}

	router := routes.SetupRoutes(customLogger, omsClient)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		stdLogger.Printf("Starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stdLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stdLogger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		stdLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	stdLogger.Println("Server exiting")
}
