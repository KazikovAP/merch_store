package model

type User struct {
	ID       int
	Username string
	Password string
	Coins    int
}

type InventoryItem struct {
	ID       int
	UserID   int
	ItemType string
	Quantity int
}

type Transaction struct {
	ID        int
	UserID    int
	Type      string
	OtherUser string
	Amount    int
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransactionDetail struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}

type InventoryResponse struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistoryResponse struct {
	Received []TransactionDetail `json:"received"`
	Sent     []TransactionDetail `json:"sent"`
}

type InfoResponse struct {
	Coins       int                 `json:"coins"`
	Inventory   []InventoryResponse `json:"inventory"`
	CoinHistory CoinHistoryResponse `json:"coinHistory"`
}
