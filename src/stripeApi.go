package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const baseApi = "https://api.stripe.com/v1"

type stripeResponse[T any] struct {
	Object  string `json:"object"`
	Data    []T    `json:"data"`
	HasMore bool   `json:"has_more"`
	Url     string `json:"url"`
}

type stripeWebhookResponse struct {
	Created int `json:"created"`
	Data    struct {
		Object checkoutSession `json:"object"`
	} `json:"data"`
}

type paymentIntent struct {
	Id      string `json:"id"`
	Amount  int    `json:"amount"`
	Created int    `json:"created"`
}

type checkoutSession struct {
	Id            string `json:"id"`
	Amount        int    `json:"amount_total"`
	Created       int    `json:"created"`
	Currency      string `json:"currency"`
	PaymentStatus string `json:"payment_status"`
	CustomFields  []struct {
		Key      string `json:"key"`
		Dropdown struct {
			Value string `json:"value"`
		} `json:"dropdown"`
		Text struct {
			Value string `json:"value"`
		} `json:"text"`
	} `json:"custom_fields"`
	CustomerDetails struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"customer_details"`
}
type donation struct {
	AmountNumber float64
	Amount       string
	Message      string
	Bootcamp     string
	Name         string
	Currency     string
}

func (cs *checkoutSession) ToDonation() donation {
	var donation donation

	for _, custom_field := range cs.CustomFields {
		if custom_field.Key == "bootcamp" {
			donation.Bootcamp = custom_field.Dropdown.Value
			continue
		}

		if custom_field.Key == "messageforthefeed" {
			donation.Message = custom_field.Text.Value
			continue
		}
	}

	donation.Amount = fmt.Sprintf("%.2f", float64(cs.Amount)/100)
	donation.AmountNumber = float64(cs.Amount) / 100

	if cs.CustomerDetails.Name != "" {
		donation.Name = strings.Fields(cs.CustomerDetails.Name)[0]
	}

	donation.Currency = cs.Currency

	return donation

}

func getPaymentIntents() []paymentIntent {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/payment_intents", baseApi), nil)

	if err != nil {
		log.Fatal("Error creating request:")
	}

	date := "2023-11-07T00:00:00Z"
	t, _ := time.Parse(time.RFC3339, date)

	q := req.URL.Query()
	q.Add("limit", "100")
	q.Add("created[gt]", fmt.Sprintf("%d", t.Unix()))

	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", stripeKey))

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:")
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Could not parse starting payment intents")
	}

	var jsonResponse stripeResponse[paymentIntent]

	json.Unmarshal(responseData, &jsonResponse)

	fmt.Printf("%v\n", jsonResponse.Data)

	return jsonResponse.Data
}

func getSessionData(id string) (checkoutSession, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/checkout/sessions?payment_intent=%s", baseApi, id), nil)

	if err != nil {
		log.Fatal("Error creating request:")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", stripeKey))

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:")
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Could not parse starting payment intents")
	}

	var jsonResponse stripeResponse[checkoutSession]

	json.Unmarshal(responseData, &jsonResponse)

	if len(jsonResponse.Data) < 1 {
		return checkoutSession{}, errors.New("no checkout found")
	}

	if jsonResponse.Data[0].PaymentStatus != "paid" {
		return checkoutSession{}, errors.New("Status not paid")
	}

	return jsonResponse.Data[0], nil
}

func getAllDonations() []donation {
	var donations []donation

	intents := getPaymentIntents()

	for _, intent := range intents {
		session, err := getSessionData(intent.Id)

		if err == nil {
			fmt.Printf("%v\n", session.ToDonation())
			donations = append(donations, session.ToDonation())
		}
	}

	return donations

	// Cant  do this calls in parallel cuz stripe gets blocked

	// var WaitGroup sync.WaitGroup
	// var sessions []checkoutSession

	// getDonationThread := func(id string, index int) {

	// 	checkoutSession, err := getSessionData(id)
	// 	if err == nil {
	// 		sessions[index] = checkoutSession
	// 	}

	// 	WaitGroup.Done()
	// }

	// intents := getPaymentIntents()

	// sessions = make([]checkoutSession, len(intents))

	// WaitGroup.Add(len(intents))

	// for index, intent := range intents {
	// 	go getDonationThread(intent.Id, index)
	// }

	// WaitGroup.Wait()

	// var output []donation

	// for _, session := range sessions {
	// 	output = append(output, session.ToDonation())
	// }
}
