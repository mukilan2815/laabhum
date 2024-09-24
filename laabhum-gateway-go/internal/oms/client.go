package oms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

    "time"

    "github.com/gin-gonic/gin"
)

type OrderType string
type OrderStatus string
type TradeStrategy string

const (
    LIMIT  OrderType = "LIMIT"
    MARKET OrderType = "MARKET"
    STOP   OrderType = "STOP"
)

const (
    PENDING   OrderStatus = "PENDING"
    EXECUTED  OrderStatus = "EXECUTED"
    CANCELLED OrderStatus = "CANCELLED"
)
const (
    SIDE_BUY           = "buy"
    SIDE_SELL          = "sell"
    TYPE_LIMIT         = "LIMIT"
    TYPE_MARKET        = "MARKET"
    TYPE_STOP          = "STOP"
    STRATEGY_SCALPING  = "scalping"
)
type Order struct {
    ID                string          `json:"id"`
    Symbol            string          `json:"symbol"`
    Quantity          int             `json:"quantity"`
    Price             float64         `json:"price"`
    Side              string          `json:"side"` // "buy" or "sell"
    Type              OrderType       `json:"type"` // LIMIT, MARKET, STOP
    Status            OrderStatus     `json:"status"` // PENDING, EXECUTED, CANCELLED
    StopPrice         float64         `json:"stop_price,omitempty"` // Stop Order Price (optional)
    Strategy          TradeStrategy   `json:"strategy"` // Trading strategy (e.g. scalping, day trading)
    RiskPercentage    float64         `json:"risk_percentage"` // % of capital risked
    StopLossActivated bool            `json:"stop_loss_activated"` // Add this field
    TakeProfit        float64         `json:"take_profit"` // Take profit level
    CreatedAt         int64           `json:"created_at"` // Timestamp for when the order is created
    ExpiresAt         time.Time       `json:"expires_at,omitempty"` // Optional expiry time for order
    ParentID          string          `json:"parent_id"` // Add ParentID field
}

type Client struct {
    BaseURL string
}

// NewClient creates a new OMS client
func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
    }
}

func (c *Client) performRequest(method, url string, body interface{}) ([]byte, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received status code %d", resp.StatusCode)
    }

    return io.ReadAll(resp.Body)
}

// StopLoss represents a stop loss structure
type StopLoss struct {
    ID      string `json:"id"`
    ParentID string `json:"parentID"`
    ChildID  string `json:"childID"`
    Active   bool   `json:"active"`
}

// ActivateStopLoss activates a stop loss for a specific order
func (c *Client) ActivateStopLoss(parentID, childID string) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/sl/%s/%s/active", c.BaseURL, parentID, childID)
    return c.performRequest(http.MethodPost, url, nil)
}

// CancelStopLoss cancels a stop loss for a specific order
func (c *Client) CancelStopLoss(parentID, childID string) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/sl/%s/%s/cancel", c.BaseURL, parentID, childID)
    return c.performRequest(http.MethodPost, url, nil)
}
func (c *Client) CreateScalperOrder(order Order) ([]byte, error) {
    log.Println("CreateScalperOrder client.go handler invoked") // Add this line for debugging
    url := fmt.Sprintf("%s/orders/scalper", c.BaseURL) // Ensure this matches your expected URL
    return c.performRequest(http.MethodPost, url, order)
}



// ExecuteOrder executes a specific order
func (c *Client) ExecuteOrder(order Order) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/%s/execute", c.BaseURL, order.ID)
    return c.performRequest(http.MethodPost, url, order)
}
// GetAllOrders retrieves all orders
func (c *Client) GetAllOrders() ([]byte, error) {
    url := fmt.Sprintf("%s/orders", c.BaseURL)
    return c.performRequest(http.MethodGet, url, nil)
}
// CancelOrder cancels a specific order
func (c *Client) CancelOrder(orderID string) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/%s/cancel", c.BaseURL, orderID)
    return c.performRequest(http.MethodDelete, url, nil)
}

// Handlers represents the HTTP handlers
type Handlers struct {
    logger    Logger
    omsClient *Client
}
// CTCOrder represents a CTC order structure
type CTCOrder struct {
    ID          string  `json:"id"`
    ParentID    string  `json:"parentID"`
    Description string  `json:"description"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
    Symbol      string  `json:"symbol"`
    OrderType   string  `json:"orderType"`
    Side        string  `json:"side"`
    ClientID    string  `json:"clientID"`
    Timestamp   string  `json:"timestamp"`
    Status      string  `json:"status"`
}

// CreateCTCOrder creates a new CTC order
func (c *Client) CreateCTCOrder(order CTCOrder) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/ctc", c.BaseURL)
    return c.performRequest(http.MethodPost, url, order)
}

// CreateCTCOrderHandler handles the creation of a new CTC order
func (h *Handlers) CreateCTCOrderHandler(c *gin.Context) {
    var order CTCOrder
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Errorf("Failed to bind CTC order data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CTC order data"})
        return
    }

    response, err := h.omsClient.CreateCTCOrder(order)
    if err != nil {
        h.logger.Errorf("Failed to create CTC order: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CTC order"})
        return
    }

    c.JSON(http.StatusOK, response)
}
// Logger is a placeholder for a logging interface
type Logger interface {
    Errorf(format string, args ...interface{})
}

// NewHandler creates a new Handler
func NewHandler(omsClient *Client, logger Logger) *Handlers {
    return &Handlers{omsClient: omsClient, logger: logger}
}

// ExecuteOrder executes a specific order
func (h *Handlers) ExecuteOrder(c *gin.Context) {
    var order Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Errorf("Failed to bind order data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
        return
    }

    response, err := h.omsClient.ExecuteOrder(order)
    if err != nil {
        h.logger.Errorf("Failed to execute order: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute order"})
        return
    }

    c.JSON(http.StatusOK, response)
}
func (h *Handlers) CreateScalperOrder(c *gin.Context) {
    var order Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Errorf("Failed to bind order data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
        return
    }

    response, err := h.omsClient.CreateScalperOrder(order)
    if err != nil {
        h.logger.Errorf("Failed to create scalper order: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scalper order"})
        return
    }

    c.JSON(http.StatusOK, response)
}

// CancelOrder cancels a specific order
func (h *Handlers) CancelOrder(c *gin.Context) {
    var order Order
    if err := c.ShouldBindJSON(&order); err != nil {
        h.logger.Errorf("Failed to bind order data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
        return
    }

    response, err := h.omsClient.CancelOrder(order.ID)
    if err != nil {
        h.logger.Errorf("Failed to cancel order: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
        return
    }

    c.JSON(http.StatusOK, response)
}

// Additional methods for Client

func (c *Client) ExecuteSpecificChild(parentID, childID string) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/%s/%s/execute", c.BaseURL, parentID, childID)
    return c.performRequest(http.MethodPost, url, nil)
}

func (c *Client) ModifyChildOrder(parentID, childID, orderType string, order Order) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/%s/%s/%s/modify", c.BaseURL, parentID, childID, orderType)
    return c.performRequest(http.MethodPatch, url, order)
}

func (c *Client) CreateOrder(order Order) ([]byte, error) {
    url := c.BaseURL + "/orders"
    body, err := json.Marshal(order)
    if err != nil {
        return nil, err
    }
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

// ExecuteAllChildTrades sends a request to execute all child trades for a parent order
func (c *Client) ExecuteAllChildTrades(parentID string) ([]byte, error) {
    url := c.BaseURL + "/orders/" + parentID + "/execute" // Ensure this is correct
    resp, err := http.Post(url, "application/json", nil)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received status code %d", resp.StatusCode)
    }
    log.Printf("Executing child trades for parentID: %s", parentID)
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) CreateCTC(parentID string, ctcOrder CTCOrder) ([]byte, error) {
    url := c.BaseURL + "/orders/" + parentID + "/ctc"
    return c.performRequest(http.MethodPost, url, ctcOrder)
}

func (c *Client) ModifyOrder(parentID, orderType string, order Order) ([]byte, error) {
    url := c.BaseURL + "/orders/" + parentID + "/" + orderType + "/modify"
    body, err := json.Marshal(order)
    if err != nil {
        return nil, err
    }
    req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) ExitAllTrades() ([]byte, error) {
    url := c.BaseURL + "/orders/exit/all"
    resp, err := http.Post(url, "application/json", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) ExitChildTrades(parentID string) ([]byte, error) {
    url := c.BaseURL + "/orders/" + parentID + "/exit"
    resp, err := http.Post(url, "application/json", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) CancelAllChildOrders(parentID string) ([]byte, error) {
    url := c.BaseURL + "/orders/" + parentID + "/cancel"
    resp, err := http.Post(url, "application/json", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) GetTrades(parentID string) ([]byte, error) {
    url := c.BaseURL + "/trades?parentID=" + parentID
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) DeleteParentOrder(parentID string) ([]byte, error) {
    url := c.BaseURL + "/orders/" + parentID
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return nil, err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) SyncPositions() ([]byte, error) {
    url := c.BaseURL + "/positions"
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func (c *Client) GetOrders() ([]byte, error) {
    url := fmt.Sprintf("%s/orders", c.BaseURL)
    return c.performRequest(http.MethodGet, url, nil)
}

func (c *Client) ExitSpecificChild(parentID, childID string) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/%s/%s/exit", c.BaseURL, parentID, childID)
    return c.performRequest(http.MethodPost, url, nil)
}

func (c *Client) CancelSpecificChildOrder(parentID, childID string) ([]byte, error) {
    url := fmt.Sprintf("%s/orders/%s/%s/cancel", c.BaseURL, parentID, childID)
    return c.performRequest(http.MethodPost, url, nil)
}

// Example usage of CancelSpecificChildOrder with advanced logic and error handling
func (h *Handlers) CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
    parentID := r.URL.Query().Get("parentID")
    orderID := r.URL.Query().Get("orderID")

    if parentID == "" || orderID == "" {
        http.Error(w, "parentID and orderID are required", http.StatusBadRequest)
        return
    }

    response, err := h.omsClient.CancelSpecificChildOrder(parentID, orderID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to cancel order: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}