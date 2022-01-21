package api_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewShopifyOrderHandler(t *testing.T) {
	t.Run("router_middleware_integration", func(t *testing.T) {
		server := NewServer()
		server.ShopifyOrderService = &mock.ShopifyOrderService{
			SaveOrderFn:     func(order shopify.Order) error { return nil },
			ProcessRefundFn: func(refund shopify.Refund) error { return nil },
		}
		server.ShopifyOrderQueueService = &mock.ShopifyOrderQueueService{
			EnqueueFn:      func(item ShopifyOrderQueueItem) error { return nil },
			ListPendingFn:  func() ([]ShopifyOrderQueueItem, error) { return []ShopifyOrderQueueItem{}, nil },
			MarkCompleteFn: func(item ShopifyOrderQueueItem) error { return nil },
		}

		r := httptest.NewRequest("POST", "/shopify/orders", strings.NewReader(shopifyOrderJSON))
		h := hmac.New(sha256.New, []byte(server.Config["SHOPIFY_SHARED_SECRET"]))
		_, _ = h.Write([]byte(shopifyOrderJSON))
		computedHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("x-shopify-hmac-sha256", computedHash)
		r.Header.Set("X-Shopify-Shop-Domain", "example.myshopify.com")
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code: %v but received %v", http.StatusOK, w.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		server := NewServer()
		server.ShopifyOrderService = &mock.ShopifyOrderService{
			SaveOrderFn:     func(order shopify.Order) error { return nil },
			ProcessRefundFn: func(refund shopify.Refund) error { return nil },
		}
		server.ShopifyOrderQueueService = &mock.ShopifyOrderQueueService{
			EnqueueFn:      func(item ShopifyOrderQueueItem) error { return nil },
			ListPendingFn:  func() ([]ShopifyOrderQueueItem, error) { return []ShopifyOrderQueueItem{}, nil },
			MarkCompleteFn: func(item ShopifyOrderQueueItem) error { return nil },
		}

		r := httptest.NewRequest("POST", "/shopify/orders", strings.NewReader(shopifyOrderJSON))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-Shopify-Shop-Domain", "some-store-name")
		w := httptest.NewRecorder()
		handler := server.NewShopifyOrderHandler()
		err := handler(w, r)
		if err != nil {
			t.Fatalf("want: %d got %d", http.StatusOK, err.StatusCode)
		}
	})

	t.Run("failure_save_order_error", func(t *testing.T) {
		// func TestShopifyOrderError(t *testing.T) {
		server := NewServer()
		server.ShopifyOrderService = &mock.ShopifyOrderService{
			SaveOrderFn:     func(order shopify.Order) error { return errors.New("error") },
			ProcessRefundFn: func(refund shopify.Refund) error { return nil },
		}

		r := httptest.NewRequest("POST", "/shopify/orders", strings.NewReader(shopifyOrderJSON))
		w := httptest.NewRecorder()
		handler := server.NewShopifyOrderHandler()
		err := handler(w, r)
		if err == nil {
			t.Fatal("expected error")
		}

		want := http.StatusInternalServerError
		if err.StatusCode != want {
			t.Fatalf("want: %d got %d", want, err.StatusCode)
		}
	})
}