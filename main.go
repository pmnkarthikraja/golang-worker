package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseData struct {
	Attributes      map[string]map[string]interface{} `json:"attributes"`
	Traits          map[string]map[string]interface{} `json:"traits"`
	Event           string                            `json:"event"`
	EventType       string                            `json:"event_type"`
	AppID           string                            `json:"app_id"`
	UserID          string                            `json:"user_id"`
	MessageID       string                            `json:"message_id"`
	PageTitle       string                            `json:"page_title"`
	PageURL         string                            `json:"page_url"`
	BrowserLanguage string                            `json:"browser_language"`
	ScreenSize      string                            `json:"screen_size"`
}

func main() {
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":8080", nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	responseCh := make(chan error)

	// Send the request to the worker goroutine
	go processRequestAndSendResponse(r, "https://webhook.site/92b74c6b-16fe-4c7b-98fd-187150acb0b0", responseCh)

	err := <-responseCh
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send success statuscode to the header
	w.WriteHeader(http.StatusOK)
}

func processRequestAndSendResponse(r *http.Request, url string, responseCh chan<- error) {
	var data ResponseData
	var inputMap map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputMap); err != nil {
		responseCh <- fmt.Errorf("error decoding JSON request: %v", err)
		return
	}
	attributes, traits := extractAttributesAndTraits(inputMap)

	data = ResponseData{
		Event:           inputMap["ev"].(string),
		EventType:       inputMap["et"].(string),
		AppID:           inputMap["id"].(string),
		UserID:          inputMap["uid"].(string),
		MessageID:       inputMap["mid"].(string),
		PageTitle:       inputMap["t"].(string),
		PageURL:         inputMap["p"].(string),
		BrowserLanguage: inputMap["l"].(string),
		ScreenSize:      inputMap["sc"].(string),
		Attributes:      attributes,
		Traits:          traits,
	}

	err := sendJSONToWebhook(url, data)
	if err != nil {
		responseCh <- err
		return
	}

	// if no error, assigning nil to channel
	responseCh <- nil
}

func extractAttributesAndTraits(inputMap map[string]interface{}) (map[string]map[string]interface{}, map[string]map[string]interface{}) {
	attributes := make(map[string]map[string]interface{})
	traits := make(map[string]map[string]interface{})

	for key, _ := range inputMap {
		switch {
		case len(key) >= 4 && key[:4] == "atrk":
			index := key[4:]
			attributeKey := inputMap["atrk"+index].(string)
			attributeType := inputMap["atrt"+index].(string)
			attributeValue := inputMap["atrv"+index].(string)
			attributes[attributeKey] = map[string]interface{}{
				"value": attributeValue,
				"type":  attributeType,
			}
		case len(key) >= 5 && key[:5] == "uatrk":
			index := key[5:]
			traitKey := inputMap["uatrk"+index].(string)
			traitType := inputMap["uatrt"+index].(string)
			traitValue := inputMap["uatrv"+index].(string)
			traits[traitKey] = map[string]interface{}{
				"value": traitValue,
				"type":  traitType,
			}
		}
	}

	return attributes, traits
}

func sendJSONToWebhook(url string, data ResponseData) error {
	jsonResp, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON response: %v", err)
	}

	respWebhook, err := http.Post(url, "application/json", bytes.NewBuffer(jsonResp))
	if err != nil {
		return fmt.Errorf("error sending request to webhook: %v", err)
	}
	defer respWebhook.Body.Close()

	if respWebhook.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code from webhook: %d", respWebhook.StatusCode)
	}

	return nil
}
