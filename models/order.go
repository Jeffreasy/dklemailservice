package models

import "time"

// WFCOrder represents a Whisky for Charity order
type WFCOrder struct {
	ID               string         `json:"id"`
	CustomerName     string         `json:"customer_name"`
	CustomerEmail    string         `json:"customer_email"`
	CustomerAddress  string         `json:"customer_address"`
	CustomerCity     string         `json:"customer_city"`
	CustomerPostal   string         `json:"customer_postal"`
	CustomerCountry  string         `json:"customer_country"`
	TotalAmount      float64        `json:"total_amount"`
	Status           string         `json:"status"`
	PaymentReference string         `json:"payment_reference,omitempty"`
	Items            []WFCOrderItem `json:"items"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty"`
}

// WFCOrderItem represents an item in a WFC order
type WFCOrderItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// WFCOrderEmailData contains data for sending WFC order emails
type WFCOrderEmailData struct {
	Order      *WFCOrder `json:"order"`
	AdminEmail string    `json:"admin_email,omitempty"`
	ToAdmin    bool      `json:"to_admin"`
	SiteURL    string    `json:"site_url,omitempty"`
}

// WFCOrderRequest represents an incoming order email request
type WFCOrderRequest struct {
	OrderID       string         `json:"order_id"`
	CustomerName  string         `json:"customer_name"`
	CustomerEmail string         `json:"customer_email"`
	TotalAmount   float64        `json:"total_amount"`
	Items         []WFCOrderItem `json:"items"`
	NotifyAdmin   bool           `json:"notify_admin,omitempty"`
	TemplateType  string         `json:"template_type"`
}
