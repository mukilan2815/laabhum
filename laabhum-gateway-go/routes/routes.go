package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Mukilan-T/laabhum-gateway-go/api" // Import the api package
	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/gorilla/mux"
)

func SetupRoutes(cfg *config.Config, logger *logger.Logger, omsClient *oms.Client) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ordersData, err := omsClient.GetOrders()
			if err != nil {
				logger.Errorf("Failed to get orders: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(ordersData)	
			return
		}
		if r.Method == http.MethodPost {
			var order oms.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				logger.Errorf("Failed to decode order: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			createdOrder, err := omsClient.CreateOrder(order)
			if err != nil {
				logger.Errorf("Failed to create order: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(createdOrder)
		}
	}).Methods(http.MethodGet, http.MethodPost)

	// Adding the new route from the second snippet
	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		handler := api.CreateOrderHandler(cfg, omsClient)
		handler(w, r)
	}).Methods("POST")

	return router
}
