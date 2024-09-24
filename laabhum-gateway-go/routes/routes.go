package routes

import (
	"github.com/Mukilan-T/laabhum-gateway-go/api"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
    "github.com/gin-gonic/gin"
    "net/http"
)

func SetupRoutes(logger *logger.Logger, omsClient *oms.Client) *gin.Engine {
	router := gin.Default()
	handlers := api.NewHandlers(logger, omsClient)
router.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
})

	// Scalper Order Routes
	router.POST("/oms/scalper/order", handlers.CreateScalperOrder)
	router.POST("/oms/scalper/order/:parentID/execute", handlers.ExecuteAllChildTrades)
	router.POST("/oms/scalper/order/:parentID/:childID/execute", handlers.ExecuteSpecificChild)
	router.POST("/oms/scalper/order/:parentID/ctc", handlers.CreateCTC)
	router.PATCH("/oms/scalper/order/:orderType/:parentID/modify", handlers.ModifyOrder)
	router.PATCH("/oms/scalper/order/:orderType/:parentID/:childID/modify", handlers.ModifyChildOrder)

	// Exit Trade Routes
	router.POST("/oms/scalper/exit/trade", handlers.ExitAllTrades)
	router.POST("/oms/scalper/trade/:parentID/exit", handlers.ExitChildTrades)
	router.POST("/oms/scalper/trade/:parentID/:childID/exit", handlers.ExitSpecificChild)

    // Cancel all child orders
    router.POST("/oms/scalper/order/:parentID/child/:childID/cancel", handlers.CancelSpecificChildOrder)
    router.POST("/oms/scalper/order/:parentID/order/:orderID/cancel", handlers.CancelSpecificOrder)

    // Get trades for a specific parent order
    router.GET("/oms/scalper/trades/:parentID", handlers.GetTrades)

    // Delete a parent order
    router.DELETE("/oms/scalper/order/:parentID", handlers.DeleteParentOrder)

    // Activate and cancel stop loss for child orders
    router.PATCH("/oms/scalper/order/sl/:parentID/:childID/active", handlers.ActivateStopLoss)
    router.PATCH("/oms/scalper/order/sl/:parentID/:childID/cancel", handlers.CancelStopLoss)

    // General Order Routes
    router.GET("/oms/orders", handlers.GetOrders)
    router.PUT("/oms/order", handlers.CreateOrder)
    router.POST("/oms/order/execute", handlers.ExecuteOrder)
    router.DELETE("/oms/order/cancel", handlers.CancelOrder)

    // Position Routes
    router.GET("/oms/positions", handlers.SyncPositions)
    router.GET("/oms/position/sync", handlers.SyncPositions)
    router.PUT("/oms/position/convert", handlers.SyncPositions)
    router.POST("/oms/position/order", handlers.CreateOrder)
    router.DELETE("/oms/position/order", handlers.CancelOrder)

    return router
}