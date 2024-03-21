package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProcessRequestAndSendResponse1(t *testing.T) {
	tests := []struct {
		name             string
		requestData      string
		ExpectedResponse ResponseData
		isEqual          bool
	}{
		{
			name: "Expected response should be equal",
			requestData: `{
				"ev": "contact_form_submitted",
				"et": "form_submit",
				"id": "cl_app_id_001",
				"uid": "cl_app_id_001-uid-001",
				"mid": "cl_app_id_001-uid-001",
				"t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
				"p": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
				"l": "en-US",
				"sc": "1920 x 1080",
				"atrk1": "form_varient",
				"atrv1": "red_top",
				"atrt1": "string",
				"atrk2": "ref",
				"atrv2": "XPOWJRICW993LKJD",
				"atrt2": "string",
				"uatrk1": "name",
				"uatrv1": "iron man",
				"uatrt1": "string",
				"uatrk2": "email",
				"uatrv2": "ironman@avengers.com",
				"uatrt2": "string",
				"uatrk3": "age",
				"uatrv3": "32",
				"uatrt3": "integer"
			}`,
			ExpectedResponse: ResponseData{
				Event:           "contact_form_submitted",
				EventType:       "form_submit",
				AppID:           "cl_app_id_001",
				UserID:          "cl_app_id_001-uid-001",
				MessageID:       "cl_app_id_001-uid-001",
				PageTitle:       "Vegefoods - Free Bootstrap 4 Template by Colorlib",
				PageURL:         "http://shielded-eyrie-45679.herokuapp.com/contact-us",
				BrowserLanguage: "en-US",
				ScreenSize:      "1920 x 1080",
				Attributes: map[string]map[string]interface{}{
					"form_varient": {"value": "red_top", "type": "string"},
					"ref":          {"value": "XPOWJRICW993LKJD", "type": "string"},
				},
				Traits: map[string]map[string]interface{}{
					"name":  {"value": "iron man", "type": "string"},
					"email": {"value": "ironman@avengers.com", "type": "string"},
					"age":   {"value": "32", "type": "integer"},
				},
			},
			isEqual: true,
		},
		{
			name: "Expected response should not be equal",
			requestData: `{
				"ev": "contact_form_submitted--wrong data",
				"et": "form_submit",
				"id": "cl_app_id_001",
				"uid": "cl_app_id_001-uid-001",
				"mid": "cl_app_id_001-uid-001",
				"t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
				"p": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
				"l": "en-US",
				"sc": "1920 x 1080",
				"atrk1": "form_varient",
				"atrv1": "red_top",
				"atrt1": "string",
				"atrk2": "ref",
				"atrv2": "XPOWJRICW993LKJD",
				"atrt2": "string",
				"uatrk1": "name",
				"uatrv1": "iron man",
				"uatrt1": "string",
				"uatrk2": "email",
				"uatrv2": "ironman@avengers.com",
				"uatrt2": "string",
				"uatrk3": "age",
				"uatrv3": "32",
				"uatrt3": "integer"
			}`,
			ExpectedResponse: ResponseData{
				Event:           "contact_form_submitted",
				EventType:       "form_submit",
				AppID:           "cl_app_id_001",
				UserID:          "cl_app_id_001-uid-001",
				MessageID:       "cl_app_id_001-uid-001",
				PageTitle:       "Vegefoods - Free Bootstrap 4 Template by Colorlib",
				PageURL:         "http://shielded-eyrie-45679.herokuapp.com/contact-us",
				BrowserLanguage: "en-US",
				ScreenSize:      "1920 x 1080",
				Attributes: map[string]map[string]interface{}{
					"form_varient": {"value": "red_top", "type": "string"},
					"ref":          {"value": "XPOWJRICW993LKJD", "type": "string"},
				},
				Traits: map[string]map[string]interface{}{
					"name":  {"value": "iron man", "type": "string"},
					"email": {"value": "ironman@avengers.com", "type": "string"},
					"age":   {"value": "32", "type": "integer"},
				},
			},
			isEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(tt.requestData)))
			if err != nil {
				t.Fatal(err)
			}

			resultCh := make(chan error)
			mockWebhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var data ResponseData
				err := json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					resultCh <- fmt.Errorf("error decoding request to webhook: %v", err)
					return
				}

				got, _ := json.Marshal(data)
				want, _ := json.Marshal(tt.ExpectedResponse)

				if tt.isEqual {
					if !reflect.DeepEqual(got, want) {
						resultCh <- fmt.Errorf("expected data not equal: got %v, want %v", data, tt.ExpectedResponse)
						return
					}
				} else {
					if reflect.DeepEqual(got, want) {
						resultCh <- fmt.Errorf("expected data equal: got %v, want %v", data, tt.ExpectedResponse)
						return
					}
				}

				resultCh <- nil
			}))

			defer mockWebhook.Close()

			go processRequestAndSendResponse(req, mockWebhook.URL, resultCh)

			//checking the result channel for any errors
			select {
			case err := <-resultCh:
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
