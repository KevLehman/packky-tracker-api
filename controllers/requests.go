package controllers

import "github.com/jinzhu/gorm"

// IntentFunc defines a function which process the data to answer intent
type IntentFunc func(*gorm.DB, WebhookRequest) WebhookResponse

// Env enclosures the global DB object
type Env struct {
	DB *gorm.DB
}

// WebhookRequest represents the request received from Fulfillment API
type WebhookRequest struct {
	ResponseID  string `json:"responseId"`
	QueryResult struct {
		QueryText  string `json:"queryText"`
		Parameters struct {
			OrderNo     string `json:"OrderNo"`
			OrderStatus string `json:"OrderStatus"`
		} `json:"parameters"`
		AllRequiredParamsPresent bool   `json:"allRequiredParamsPresent"`
		FulfillmentText          string `json:"fulfillmentText"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text"`
			} `json:"text"`
			Platform string `json:"platform,omitempty"`
		} `json:"fulfillmentMessages"`
		Intent struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"intent"`
		IntentDetectionConfidence float64 `json:"intentDetectionConfidence"`
		LanguageCode              string  `json:"languageCode"`
	} `json:"queryResult"`
	OriginalDetectIntentRequest struct {
		Payload struct {
		} `json:"payload"`
	} `json:"originalDetectIntentRequest"`
	Session string `json:"session"`
}

// WebhookCard represents a Messenger Card payload
type WebhookCard struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURI string `json:"imageUri"`
	Platform string `json:"platform"`
	Lang     string `json:"lang"`
	Buttons  []struct {
		Text     string `json:"text"`
		Postback string `json:"postback"`
	} `json:"buttons"`
}

// WebhookResponse is the response we send to fulfillment API
type WebhookResponse struct {
	FulfillmentText     string `json:"fulfillmentText"`
	FulfillmentMessages []struct {
		WebhookCard `json:"card"`
	} `json:"fulfillmentMessages"`
	Source string `json:"source"`
	/* Payload struct {
		Facebook struct {
			Text string `json:"text"`
		} `json:"facebook"`
		Slack struct {
			Text string `json:"text"`
		} `json:"slack"`
	} `json:"payload"` */
	/* OutputContexts []struct {
		Name          string `json:"name"`
		LifespanCount int    `json:"lifespanCount"`
		Parameters    struct {
			Param string `json:"param"`
		} `json:"parameters"`
	} `json:"outputContexts"`
	FollowupEventInput struct {
		Name         string `json:"name"`
		LanguageCode string `json:"languageCode"`
		Parameters   struct {
			Param string `json:"param"`
		} `json:"parameters"`
	} `json:"followupEventInput"` */
}
