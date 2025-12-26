package handler

import (
	"miniapp-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentSvc *service.PaymentService
}

func NewPaymentHandler(paymentSvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

type CreateOrderRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
	UserID int64  `json:"user_id"` // Optional if using middleware, but here we might need it
	OpenID string `json:"openid"`  // Usually from context or DB
}

func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, get UserID from JWT token in context
	// For now, we assume it's passed or we get it from query/header if not in body
	// Let's assume the frontend passes user_id for simplicity as per previous patterns
	if req.UserID == 0 {
		// Try to get from query if not in JSON (though it should be in JSON for POST)
		// Or return error
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// We need OpenID for JSAPI payment. 
	// We can fetch it from DB using UserID if not provided.
	// Let's assume we need to fetch it.
	// Since I don't have direct access to repo here, I should probably pass UserID to service 
	// and let service fetch OpenID. 
	// But my service method signature is CreateOrder(userID, openID, ...).
	// I'll update the service to fetch OpenID if I can, or I'll fetch it here if I had the repo.
	// Wait, PaymentService has UserRepo. I should update PaymentService to fetch OpenID.
	
	// For now, let's assume the service handles it or we pass it.
	// Actually, the service method I wrote takes `openID`. 
	// I should update the service to look up the user by ID to get the OpenID.
	
	// Let's update the handler to just pass UserID and let Service handle OpenID lookup.
	// But wait, I already wrote the service to take OpenID.
	// I will modify the service in the next step to look up OpenID.
	
	// For now, let's assume the frontend sends it or I'll fix the service.
	// Actually, I'll fix the service.
	
	clientIP := c.ClientIP()
	
	// Call service (assuming I'll update it to not need explicit OpenID or I'll fetch it)
	// Let's pass empty OpenID and let service fetch it.
	
	params, err := h.paymentSvc.CreateOrder(req.UserID, req.OpenID, req.PlanID, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, params)
}

func (h *PaymentHandler) Notify(c *gin.Context) {
	// Pass the request to the service to handle verification and business logic
	// The SDK expects standard http.Request and ResponseWriter
	err := h.paymentSvc.HandleNotify(c.Request, c.Writer)
	if err != nil {
		// If verification fails or other error, return error
		// Note: SDK might have already written response if we used its handler wrapper, 
		// but here we are calling VerifySign manually.
		// We need to return XML response to WeChat.
		// Actually, s.pay.VerifySign just returns the parsed data.
		// We need to reply to WeChat.
		c.String(http.StatusInternalServerError, "FAIL")
		return
	}

	// Reply success to WeChat
	c.String(http.StatusOK, `<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>`)
}

func (h *PaymentHandler) GetSubscription(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	isVIP, expireDate, err := h.paymentSvc.GetSubscription(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_vip":      isVIP,
		"expire_date": expireDate,
	})
}
