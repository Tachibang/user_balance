package entity

import "time"

type Operation struct {
	Id            int        `json:"id"`
	AccountId     int        `json:"account_id"`
	Amount        int        `json:"amount"`
	OperationType string     `json:"operation_type"`
	ProductId     *int       `json:"product_id,omitempty"`
	Description   *string    `json:"description,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
