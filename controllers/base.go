package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var configuredIntents = map[string]IntentFunc{
	"Track Order":          TrackOrderIntent,
	"Start Order Tracking": StartTrackingOnProduct,
	"Update Order Status":  UpdateOrderStatus,
}

// IntentController defines the main controller for Dialogflow intents
func (env *Env) IntentController(c *gin.Context) {
	fulfillmentResponse := WebhookRequest{}
	if err := c.ShouldBindJSON(&fulfillmentResponse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	intent := fulfillmentResponse.QueryResult.Intent.DisplayName
	intentFunc, ok := configuredIntents[intent]

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Provided intent has no matching function",
		})
		return
	}

	log.Println("Matched intent: ", intent)
	c.JSON(http.StatusOK, intentFunc(env.DB, fulfillmentResponse))
	return
}
