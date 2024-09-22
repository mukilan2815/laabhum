package oms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Order struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Symbol      string  `json:"symbol"`
	OrderType   string  `json:"orderType"`
	Side        string  `json:"side"`
	ClientID    string  `json:"clientID"`
	Timestamp   string  `json:"timestamp"`
}

type Client struct {
	BaseURL string
	// Add necessary fields
}

// NewClient creates a new OMS client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

// CreateOrder creates a new order
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
	return ioutil.ReadAll(resp.Body)
}

// ExecuteChildOrder executes a child order based on parentID and childID
func (c *Client) ExecuteChildOrder(parentID, childID string) ([]byte, error) {
	url := c.BaseURL + "/orders/execute"
	reqBody := map[string]string{"parentID": parentID, "childID": childID}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetTrades retrieves trades based on parentID
func (c *Client) GetTrades(parentID string) ([]byte, error) {
	url := c.BaseURL + "/trades?parentID=" + parentID
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetOrders retrieves all orders
func (c *Client) GetOrders() ([]byte, error) {
	url := c.BaseURL + "/orders"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// OMSClient is a new client structure
type OMSClient struct {
	baseURL string
}

// NewOMSClient creates a new OMSClient
func NewOMSClient() *OMSClient {
	return &OMSClient{
		baseURL: "http://localhost:8080", // URL for OMS
	}
}

// CreateOrder creates a new order using OMSClient
func (c *OMSClient) CreateOrder(order map[string]interface{}) error {
	orderPayload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/orders", bytes.NewBuffer(orderPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create order")
	}

	return nil
}
