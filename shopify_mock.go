package mock

import (
	"<path to package>/shopify"
)

// ShopifyOrderService represents a mock implementation of shopify.OrderService.
type ShopifyOrderService struct {
	SaveOrderFn          func(order shopify.Order) error
	SaveOrderInvoked     bool
	SaveOrderInvokedWith shopify.Order

	ProcessRefundFn          func(refund shopify.Refund) error
	ProcessRefundInvoked     bool
	ProcessRefundInvokedWith shopify.Refund
}

// SaveOrder invokes the mock implementation and marks the function as invoked.
func (s *ShopifyOrderService) SaveOrder(order shopify.Order) error {
	s.SaveOrderInvoked = true
	s.SaveOrderInvokedWith = order
	return s.SaveOrderFn(order)
}

// ProcessRefund invokes the mock implementation and marks the function as invoked.
func (s *ShopifyOrderService) ProcessRefund(refund shopify.Refund) error {
	s.ProcessRefundInvoked = true
	s.ProcessRefundInvokedWith = refund
	return s.ProcessRefundFn(refund)
}

func (s *ShopifyOrderService) Get(ID string) (shopify.Order, error) {
	return shopify.Order{}, nil
}