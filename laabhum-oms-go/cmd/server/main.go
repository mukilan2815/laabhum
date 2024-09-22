package main

import (
	"encoding/json"
	"log"
	"net/http"

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

	log.Println("Server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func ordersHandler(omsService *service.OMSService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			filter := repository.OrderFilter{}
			orders, err := omsService.GetOrders(filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response, err := json.Marshal(orders)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else if r.Method == http.MethodPost {
			var order models.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			createdOrder, err := omsService.CreateOrder(order)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response, err := json.Marshal(createdOrder)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		}
	}
}
