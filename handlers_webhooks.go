package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"<path to package>/shopify"
)

// This will be called by the Shopify fulfillments CREATE webhook
// which happens on every completed shipment. 
func (s *Server) NewShopifyFulfillmentsHandler() handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fulfillment := shopify.Fulfillment{}

		// Duplicate body for payload
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.NewDecoder(r.Body).Decode(&fulfillment)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shopDomain := r.Header.Get("X-Shopify-Shop-Domain")
		if len(shopDomain) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jobContext := "shopify"

		// Save webhook payload context and payload, to be processed later
		job := WebhookJob{Context: jobContext, Payload: string(body)}
		err = s.WebhookJobService.Create(job)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return nil
	}
}

// NewShopifyOrder is called by the Shopify order COMPLETE webhhook, which
// means all shipments have been completed and the order is closed.
func (s *Server) NewShopifyOrderHandler() handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shopDomain := r.Header.Get("X-Shopify-Shop-Domain")
		if len(shopDomain) < 1 {
			return ErrInternalServer(fmt.Errorf("No X-Shopify-Shop-Domain header in request headers: %v", r.Header))
		}

		order := shopify.Order{}
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.ShopifyOrderService.SaveOrder(order)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return nil
	}
}

// NewShopifyRefund will be called via webhook from Shopify
func (s *Server) NewShopifyRefundHandler() handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refund := shopify.Refund{}
		err := json.NewDecoder(r.Body).Decode(&refund)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = s.ShopifyOrderService.ProcessRefund(refund)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return nil
	}
}