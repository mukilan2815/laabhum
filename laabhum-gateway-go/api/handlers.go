package api

import (
	"fmt"
	"net/http"
	"time"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger" // Add this line
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	logger    *logger.Logger
	omsClient *oms.Client
}

func NewHandlers(logger *logger.Logger, omsClient *oms.Client) *Handlers {
	return &Handlers{
		logger:    logger,
		omsClient: omsClient,
	}
}
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
// Error handler utility function
func (h *Handlers) handleError(c *gin.Context, statusCode int, err error, msg string) {
	h.logger.Errorf("%s: %v", msg, err)
	c.JSON(statusCode, gin.H{"error": msg})
}

// CreateOrder handles creating a new order
func (h *Handlers) CreateOrder(c *gin.Context) {
	var order oms.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		h.handleError(c, http.StatusBadRequest, err, "Invalid order data")
		return
	}

	response, err := h.omsClient.CreateOrder(order)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to create order")
		return
	}

	c.JSON(http.StatusCreated, response)
}
// CancelSpecificOrder cancels a specific order
func (h *Handlers) CancelSpecificOrder(c *gin.Context) {
    orderID := c.Param("orderID")
    if orderID == "" {
        h.handleError(c, http.StatusBadRequest, nil, "Order ID is required")
        return
    }

    response, err := h.omsClient.CancelOrder(orderID)
    if err != nil {
        h.handleError(c, http.StatusInternalServerError, err, "Failed to cancel order")
        return
    }

    c.JSON(http.StatusOK, response)
}
// CancelOrder cancels an order
func (h *Handlers) CancelOrder(c *gin.Context) {
	orderID := c.Param("orderID")
	if orderID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Order ID is required")
		return
	}

	response, err := h.omsClient.CancelOrder(orderID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to cancel order")
		return
	}

	c.JSON(http.StatusOK, response)
}
// ActivateStopLoss activates a stop loss for a parent order
func (h *Handlers) ActivateStopLoss(c *gin.Context) {
    parentID := c.Param("parentID")
    if parentID == "" {
        h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
        return
    }

    var stopLoss oms.StopLoss
    if err := c.ShouldBindJSON(&stopLoss); err != nil {
        h.handleError(c, http.StatusBadRequest, err, "Invalid stop loss data")
        return
    }

	response, err := h.omsClient.ActivateStopLoss(parentID, stopLoss.ID) // Assuming stopLoss has an ID field of type string
    if err != nil {
        h.handleError(c, http.StatusInternalServerError, err, "Failed to activate stop loss")
        return
    }

    c.JSON(http.StatusOK, response)
}

// CancelStopLoss cancels a stop loss for a parent order
func (h *Handlers) CancelStopLoss(c *gin.Context) {
    parentID := c.Param("parentID")
    if parentID == "" {
        h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
        return
    }

    childID := c.Param("childID")
    if childID == "" {
        h.handleError(c, http.StatusBadRequest, nil, "Child ID is required")
        return
    }

	response, err := h.omsClient.CancelStopLoss(parentID, childID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to cancel stop loss")
		return
	}

	c.JSON(http.StatusOK, response)
}
func (h *Handlers) CreateScalperOrder(c *gin.Context) {
	var order oms.Order // Assuming ScalperOrder should be of type Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Printf("Invalid input for scalper order: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    // Log received data for debugging
    h.logger.Printf("Received order data: %+v\n", order)

    // Validate required fields
    if order.Quantity <= 0 || order.Price <= 0 || order.Symbol == "" || 
       order.Type == "" || order.Side == "" || 
       order.Strategy == "" || 
       order.RiskPercentage <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
        return
    }

    // Create the scalper order via service layer
	createdOrder, err := h.omsClient.CreateScalperOrder(order)
    if err != nil {
        h.logger.Printf("Order creation failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Order creation failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Scalper order created successfully", "order": createdOrder})
}
// validateOrder checks the order fields for validity
func validateOrder(order oms.Order) error {
    if order.Quantity <= 0 || order.Price <= 0 || order.Symbol == "" || order.Type == "" || order.Side == "" || order.Strategy == "" || order.RiskPercentage <= 0 {
        return fmt.Errorf("Invalid order data")
    }
    
    // Optional validation for stop loss
    if order.StopLossActivated && order.StopPrice <= 0 {
        return fmt.Errorf("Stop price must be greater than 0 if stop loss is activated")
    }
    
    return nil
}

func (h *Handlers) ExecuteAllChildTrades(c *gin.Context) {
    parentID := c.Param("parentID")
    if parentID == "" {
        h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
        return
    }

    response, err := h.omsClient.ExecuteAllChildTrades(parentID)
    if err != nil {
        h.handleError(c, http.StatusInternalServerError, err, "Failed to execute all child trades")
        return
    }

    c.JSON(http.StatusOK, response)
}

// ExecuteSpecificChild executes a specific child trade
func (h *Handlers) ExecuteSpecificChild(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")
	if parentID == "" || childID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID and Child ID are required")
		return
	}

	response, err := h.omsClient.ExecuteSpecificChild(parentID, childID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to execute specific child")
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateCTC creates a CTC for a parent order
func (h *Handlers) CreateCTC(c *gin.Context) {
	parentID := c.Param("parentID")
	if parentID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
		return
	}

	var ctcOrder oms.CTCOrder
	if err := c.ShouldBindJSON(&ctcOrder); err != nil {
		h.handleError(c, http.StatusBadRequest, err, "Invalid CTC order data")
		return
	}

	response, err := h.omsClient.CreateCTC(parentID, ctcOrder)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to create CTC")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ModifyOrder modifies an existing order based on the type
func (h *Handlers) ModifyOrder(c *gin.Context) {
	parentID := c.Param("parentID")
	orderType := c.Param("orderType")
	if parentID == "" || orderType == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID and Order Type are required")
		return
	}

	var order oms.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		h.handleError(c, http.StatusBadRequest, err, "Invalid order data")
		return
	}

	response, err := h.omsClient.ModifyOrder(parentID, orderType, order)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to modify order")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ModifyChildOrder modifies an existing child order
func (h *Handlers) ModifyChildOrder(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")
	orderType := c.Param("orderType")
	if parentID == "" || childID == "" || orderType == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID, Child ID, and Order Type are required")
		return
	}

	var order oms.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		h.handleError(c, http.StatusBadRequest, err, "Invalid order data")
		return
	}

	response, err := h.omsClient.ModifyChildOrder(parentID, childID, orderType, order)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to modify child order")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ExitAllTrades exits all trades for scalper
func (h *Handlers) ExitAllTrades(c *gin.Context) {
	response, err := h.omsClient.ExitAllTrades()
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to exit all trades")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ExitChildTrades exits all child trades for a parent order
func (h *Handlers) ExitChildTrades(c *gin.Context) {
	parentID := c.Param("parentID")
	if parentID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
		return
	}

	response, err := h.omsClient.ExitChildTrades(parentID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to exit child trades")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ExitSpecificChild exits a specific child trade
func (h *Handlers) ExitSpecificChild(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")
	if parentID == "" || childID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID and Child ID are required")
		return
	}

	response, err := h.omsClient.ExitSpecificChild(parentID, childID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to exit specific child")
		return
	}

	c.JSON(http.StatusOK, response)
}

// CancelAllChildOrders cancels all child orders for a parent order
func (h *Handlers) CancelAllChildOrders(c *gin.Context) {
	parentID := c.Param("parentID")
	if parentID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
		return
	}

	response, err := h.omsClient.CancelAllChildOrders(parentID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to cancel all child orders")
		return
	}

	c.JSON(http.StatusOK, response)
}

// CancelSpecificChildOrder cancels a specific child order
func (h *Handlers) CancelSpecificChildOrder(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")
	if parentID == "" || childID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID and Child ID are required")
		return
	}

	response, err := h.omsClient.CancelSpecificChildOrder(parentID, childID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to cancel specific child order")
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTrades retrieves trades based on parentID
func (h *Handlers) GetTrades(c *gin.Context) {
	parentID := c.Param("parentID")
	if parentID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
		return
	}

	response, err := h.omsClient.GetTrades(parentID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to get trades")
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteParentOrder deletes a parent order
func (h *Handlers) DeleteParentOrder(c *gin.Context) {
	parentID := c.Param("parentID")
	if parentID == "" {
		h.handleError(c, http.StatusBadRequest, nil, "Parent ID is required")
		return
	}

	response, err := h.omsClient.DeleteParentOrder(parentID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to delete parent order")
		return
	}

	c.JSON(http.StatusOK, response)
}

// SyncPositions syncs the current positions
func (h *Handlers) SyncPositions(c *gin.Context) {
	response, err := h.omsClient.SyncPositions()
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to sync positions")
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetOrders retrieves all orders from the OMS
func (h *Handlers) GetOrders(c *gin.Context) {
	response, err := h.omsClient.GetAllOrders()
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to retrieve orders")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ExecuteOrder executes an order
func (h *Handlers) ExecuteOrder(c *gin.Context) {
	var order oms.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		h.handleError(c, http.StatusBadRequest, err, "Invalid order data")
		return
	}

	response, err := h.omsClient.ExecuteOrder(order)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Failed to execute order")
		return
	}

	c.JSON(http.StatusOK, response)
}