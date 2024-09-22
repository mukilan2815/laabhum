package api

import (
	"encoding/json"
	"net/http"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/Mukilan-T/laabhum-oms-go/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type Handlers struct {
	omsService *service.OMSService
}

// NewHandlers initializes the handlers with OMSService
func NewHandlers(omsService *service.OMSService) *Handlers {
	return &Handlers{
		omsService: omsService,
	}
}

// CreateScalperOrder creates a new scalper order
func (h *Handlers) CreateScalperOrder(c *gin.Context) {
	var order models.ScalperOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	createdOrder, err := h.omsService.CreateScalperOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order creation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Scalper order created successfully", "order": createdOrder})
}

// ExecuteChildOrder executes a specific child order
func (h *Handlers) ExecuteChildOrder(c *gin.Context) {
	// parentID := c.Param("parentID")
	childID := c.Param("childID")

	err := h.omsService.ExecuteChildOrder(childID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute child order: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Child order executed successfully"})
}

// GetTrades retrieves trades based on parent order ID
func (h *Handlers) GetTrades(c *gin.Context) {
	parentID := c.Param("parentId")

	trades, err := h.omsService.GetTrades(parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve trades: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// CreateOrder creates a new order
func (h *Handlers) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	createdOrder, err := h.omsService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order creation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order": createdOrder})
}

// GetOrders retrieves all orders with optional filters
func (h *Handlers) GetOrders(c *gin.Context) {
	orders, err := h.omsService.GetOrders(repository.OrderFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// SetupRoutes initializes all the routes
func SetupRoutes(repo repository.OrderRepository, omsService service.OMSService) *mux.Router {
	router := mux.NewRouter()

	// Route to create and get orders
	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var order models.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
				return
			}

			createdOrder, err := omsService.CreateOrder(order)
			if err != nil {
				http.Error(w, "Order creation failed: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Order created successfully",
				"order":   createdOrder,
			})
		case http.MethodGet:
			orders, err := omsService.GetOrders(repository.OrderFilter{})
			if err != nil {
				http.Error(w, "Failed to retrieve orders: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(orders)
		}
	}).Methods(http.MethodGet, http.MethodPost)

	return router
}

// CreateOrderHandler is a simpler handler for order processing
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call OMS service to process the order
	err = service.ProcessOrder(order)
	if err != nil {
		http.Error(w, "Failed to process order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order processed successfully"})
}
