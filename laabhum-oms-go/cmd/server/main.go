package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/Mukilan-T/laabhum-oms-go/service"
	"github.com/gorilla/mux"
)

func main() {
	repo := repository.NewInMemoryOrderRepository()
	omsService := service.NewOMSService(repo)

	r := mux.NewRouter()
	r.HandleFunc("/orders", ordersHandler(omsService)).Methods(http.MethodGet, http.MethodPost)

	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		log.Println("Server started at :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Fatalf("Server close error: %v", err)
	}
}

func ordersHandler(omsService *service.OMSService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			filter := repository.OrderFilter{}
			orders, err := omsService.GetOrders(filter)
			if err != nil {
				http.Error(w, "Error retrieving orders: "+err.Error(), http.StatusInternalServerError)
				return
			}
			response, err := json.Marshal(orders)
			if err != nil {
				http.Error(w, "Error marshalling response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else if r.Method == http.MethodPost {
			var order models.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, "Invalid order data: "+err.Error(), http.StatusBadRequest)
				return
			}
			createdOrder, err := omsService.CreateOrder(order)
			if err != nil {
				http.Error(w, "Error creating order: "+err.Error(), http.StatusInternalServerError)
				return
			}
			response, err := json.Marshal(createdOrder)
			if err != nil {
				http.Error(w, "Error marshalling response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}
	}
}
