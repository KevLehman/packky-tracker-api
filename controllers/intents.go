package controllers

import (
	"fmt"
	"log"

	"github.com/KevLehmann/packky-tracker-api/models"
	"github.com/jinzhu/gorm"
)

// TrackOrderIntent evaluates payload and returns the necessary data to Dialogflow
func TrackOrderIntent(db *gorm.DB, payload WebhookRequest) WebhookResponse {
	response := WebhookResponse{}
	order := payload.QueryResult.Parameters.OrderNo

	foundOrder, err := models.Item{}.GetItemByOrder(db, order)

	if err != nil {
		log.Println("[TrackOrderIntent] Error while finding Order: ", err.Error())
		response.FulfillmentText = fmt.Sprintf("El pedido %s no existe en nuestros registros", order)
		return response
	}

	response.FulfillmentText = fmt.Sprintf("Orden %s | Estado: %s", order, foundOrder.Status.Name)
	response.FulfillmentMessages = []struct {
		WebhookCard `json:"card"`
	}{
		{
			WebhookCard{
				Title:    fmt.Sprintf("Pedido %s", order),
				Subtitle: fmt.Sprintf("Estado del pedido: %s", foundOrder.Status.Name),
				ImageURI: foundOrder.Status.ImageURI,
				Platform: "FACEBOOK",
				Lang:     "es",
			},
		},
	}

	response.Source = "Packky-API"
	return response
}

// StartTrackingOnProduct sets tracking status capabilities to a product
func StartTrackingOnProduct(db *gorm.DB, payload WebhookRequest) WebhookResponse {
	response := WebhookResponse{
		Source: "Packky-API",
	}
	order := payload.QueryResult.Parameters.OrderNo
	initialState, err := models.ItemStatus{}.GetInitialStatus(db)

	if err != nil {
		log.Println("[StartTrackingOnProduct] Error while finding State: ", err.Error())
		response.FulfillmentText = fmt.Sprintf("No podemos encontrar un estado asignable a la orden especificada\n")
		return response
	}

	foundOrder, err := models.Item{}.FindOrCreate(db, order)
	if err != nil {
		log.Println("[StartTrackingOnProduct] Error while finding Order: ", err.Error())
		response.FulfillmentText = fmt.Sprintf("Algo ocurrio mientras procesabamos la orden, intenta de nuevo en un momento\n")
		return response
	}
	foundOrder.Status = *initialState
	foundOrder.Update(db)

	response.FulfillmentText = fmt.Sprintf("Pedido %s | Estado %s\n", order, initialState.Name)

	return response
}

// UpdateOrderStatus takes a new status for an order and moves the product to it
func UpdateOrderStatus(db *gorm.DB, payload WebhookRequest) WebhookResponse {
	response := WebhookResponse{
		Source: "Packky-API",
	}
	order := payload.QueryResult.Parameters.OrderNo
	newStatus := payload.QueryResult.Parameters.OrderStatus

	foundOrder, err := models.Item{}.GetItemByOrder(db, order)
	if err != nil {
		log.Println("[UpdateOrderStatus] Error while finding Order: ", err.Error())
		response.FulfillmentText = "Algo ocurrio mientras buscabamos la orden, intenta mas tarde"
		return response
	}
	foundStatus, err := models.ItemStatus{}.GetStatusByName(db, newStatus)
	if err != nil {
		log.Println("[UpdateOrderStatus] Error while finding Status: ", err.Error())
		response.FulfillmentText = "Algo ocurrio mientras buscabamos la orden, intenta mas tarde"
		return response
	}

	foundOrder.Status = *foundStatus

	if err := foundOrder.Update(db); err != nil {
		response.FulfillmentText = "No pudimos actualizar el estado de tu orden, intenta de nuevo mas tarde"
		log.Println("[UpdateOrderStatus] Error while updating Order: ", err.Error())
		return response
	}

	response.FulfillmentText = fmt.Sprintf("Pedido %s | Estado %s\n", order, newStatus)
	return response
}
