package api

import (
	"log"
	"net/http"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/Mukilan-T/laabhum-oms-go/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
    logger     *log.Logger
    omsService *service.OMSService
}

// NewHandlers initializes the handlers with OMSService
func NewHandlers(logger *log.Logger, omsService *service.OMSService) *Handlers {
    return &Handlers{
        logger:     logger,
        omsService: omsService,
    }
}
func SetupRoutes(logger *log.Logger, omsService *service.OMSService) *gin.Engine {
	router := gin.Default()
	handlers := NewHandlers(logger, omsService)

	// Scalper Order Routes
	router.POST("/oms/scalper/order", handlers.CreateScalperOrder)
	router.POST("/oms/scalper/order/:parentID/execute", handlers.ExecuteAllChildTrades)
	router.POST("/oms/scalper/order/:parentID/:childID/execute", handlers.ExecuteSpecificChild)
	router.POST("/oms/scalper/order/:parentID/ctc", handlers.CreateCTC)
	router.PATCH("/oms/scalper/order/:orderType/:parentID/modify", handlers.ModifyOrder)
	router.PATCH("/oms/scalper/order/:orderType/:parentID/:childID/modify", handlers.ModifyChildOrder)

	// Exit scalper trades
	router.POST("/oms/scalper/exit/trade", handlers.ExitAllTrades)
	router.POST("/oms/scalper/trade/:parentID/exit", handlers.ExitChildTrades)
	router.POST("/oms/scalper/trade/:parentID/:childID/exit", handlers.ExitSpecificChild)

	// Cancel scalper orders
	router.POST("/oms/scalper/order/:parentID/cancel", handlers.CancelAllChildOrders)
	router.POST("/oms/scalper/order/:parentID/:childID/cancel", handlers.CancelSpecificChildOrder)

	// Get trades for a specific parent order
	router.GET("/oms/scalper/trades/:parentID", handlers.GetTrades)

	// Delete a parent order
	router.DELETE("/oms/scalper/order/:parentID", handlers.DeleteParentOrder)


	// Position Routes
	router.GET("/oms/positions", handlers.SyncPositions)

	// General Order Routes
	router.GET("/oms/orders", handlers.GetOrders)
	router.PUT("/oms/order", handlers.CreateOrder)
	router.POST("/oms/order/execute", handlers.ExecuteOrder)
	router.DELETE("/oms/order/cancel", handlers.CancelOrder)

	return router
}

// CreateScalperOrder handles creating a scalper order
func (h *Handlers) CreateScalperOrder(c *gin.Context) {
    var order models.ScalperOrder
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Printf("Invalid input for scalper order: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }
    h.logger.Println("CreateScalperOrder handlers.go in oms handler invoked") // Add this line for debugging

    createdOrder, err := h.omsService.CreateScalperOrder(order)
    if err != nil {
        h.logger.Printf("Order creation failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Order creation failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Scalper order created successfully", "order": createdOrder})
}

// ExecuteAllChildTrades executes all child trades for a parent order
func (h *Handlers) ExecuteAllChildTrades(c *gin.Context) {
    parentID := c.Param("parentID")
    
    err := h.omsService.ExecuteAllChildTrades(parentID)
    if err != nil {
        h.logger.Printf("Failed to execute all child trades: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute all child trades: " + err.Error()})
        return
    }

    h.logger.Printf("Executed all child trades for parentID: %s", parentID)
    c.JSON(http.StatusOK, gin.H{"message": "All child trades executed successfully"})
}

// ExecuteSpecificChild executes a specific child trade
func (h *Handlers) ExecuteSpecificChild(c *gin.Context) {
    parentID := c.Param("parentID")
    childID := c.Param("childID")

    err := h.omsService.ExecuteSpecificChild(parentID, childID)
    if err != nil {
        h.logger.Printf("Failed to execute specific child trade: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute specific child trade: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Specific child trade executed successfully"})
}

// CreateCTC creates a CTC order
func (h *Handlers) CreateCTC(c *gin.Context) {
    var ctcOrder models.CTCOrder
    if err := c.ShouldBindJSON(&ctcOrder); err != nil {
        h.logger.Printf("Invalid input for CTC order: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    createdOrder, err := h.omsService.CreateCTC(ctcOrder)
    if err != nil {
        h.logger.Printf("CTC order creation failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "CTC order creation failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "CTC order created successfully", "order": createdOrder})
}

// ExitAllTrades exits all trades
func (h *Handlers) ExitAllTrades(c *gin.Context) {
    err := h.omsService.ExitAllTrades("someStringArgument")
    if err != nil {
        h.logger.Printf("Failed to exit all trades: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exit all trades: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "All trades exited successfully"})
}

// ExitChildTrades exits child trades for a given parent order
func (h *Handlers) ExitChildTrades(c *gin.Context) {
    parentID := c.Param("parentID")

    err := h.omsService.ExitChildTrades(parentID)
    if err != nil {
        h.logger.Printf("Failed to exit child trades: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exit child trades: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Child trades exited successfully"})
}

// ExitSpecificChild exits a specific child trade
func (h *Handlers) ExitSpecificChild(c *gin.Context) {
    parentID := c.Param("parentID")
    childID := c.Param("childID")

    err := h.omsService.ExitSpecificChild(parentID, childID)
    if err != nil {
        h.logger.Printf("Failed to exit specific child trade: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exit specific child trade: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Specific child trade exited successfully"})
}

// CancelAllChildOrders cancels all child orders for a given parent order
func (h *Handlers) CancelAllChildOrders(c *gin.Context) {
    parentID := c.Param("parentID")

    err := h.omsService.CancelAllChildOrders(parentID)
    if err != nil {
        h.logger.Printf("Failed to cancel all child orders: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel all child orders: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "All child orders cancelled successfully"})
}

// CancelSpecificChildOrder cancels a specific child order
func (h *Handlers) CancelSpecificChildOrder(c *gin.Context) {
    parentID := c.Param("parentID")
    childID := c.Param("childID")

    err := h.omsService.CancelSpecificChildOrder(parentID, childID)
    if err != nil {
        h.logger.Printf("Failed to cancel specific child order: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel specific child order: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Specific child order cancelled successfully"})
}

// DeleteParentOrder deletes a parent order
func (h *Handlers) DeleteParentOrder(c *gin.Context) {
    parentID := c.Param("parentID")

    err := h.omsService.DeleteParentOrder(parentID)
    if err != nil {
        h.logger.Printf("Failed to delete parent order: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete parent order: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Parent order deleted successfully"})
}

// ActivateStopLoss activates stop loss for a specific child order

// SyncPositions syncs positions
func (h *Handlers) SyncPositions(c *gin.Context) {
    err := h.omsService.SyncPositions()
    if err != nil {
        h.logger.Printf("Failed to sync positions: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync positions: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Positions synced successfully"})
}

// ExecuteOrder executes an order
func (h *Handlers) ExecuteOrder(c *gin.Context) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Printf("Invalid input for order execution: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    err := h.omsService.ExecuteOrder(order)
    if err != nil {
        h.logger.Printf("Order execution failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Order execution failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Order executed successfully"})
}

// CancelOrder cancels an order
func (h *Handlers) CancelOrder(c *gin.Context) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Printf("Invalid input for order cancellation: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    err := h.omsService.CancelOrder(order.ID) // Assuming 'ID' is the string field required
    if err != nil {
        h.logger.Printf("Order cancellation failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Order cancellation failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// GetTrades retrieves trades for a given parent ID
func (h *Handlers) GetTrades(c *gin.Context) {
    parentID := c.Param("parentID")

    trades, err := h.omsService.GetTrades(parentID)
    if err != nil {
        h.logger.Printf("Failed to retrieve trades: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve trades: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, trades)
}

// CreateOrder handles creating a new order
func (h *Handlers) CreateOrder(c *gin.Context) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Printf("Invalid input for order: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    createdOrder, err := h.omsService.CreateOrder(order)
    if err != nil {
        h.logger.Printf("Order creation failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Order creation failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully", "order": createdOrder})
}

// ModifyOrder modifies an existing order
func (h *Handlers) ModifyOrder(c *gin.Context) {
    // Implement logic to modify an order here
    c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// ModifyChildOrder modifies an existing child order
func (h *Handlers) ModifyChildOrder(c *gin.Context) {
    // Implement logic to modify a child order here
    c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetOrders retrieves all orders
func (h *Handlers) GetOrders(c *gin.Context) {
    orders, err := h.omsService.GetOrders(repository.OrderFilter{})
    if err != nil {
        h.logger.Printf("Failed to retrieve orders: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, orders)
}