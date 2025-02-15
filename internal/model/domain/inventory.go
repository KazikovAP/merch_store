package domain

type InventoryItem struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userId"`
	ItemType string `json:"itemType"`
	Quantity int    `json:"quantity"`
}
