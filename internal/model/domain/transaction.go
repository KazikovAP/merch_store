package domain

type Transaction struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Type      string `json:"type"`
	OtherUser string `json:"otherUser"`
	Amount    int    `json:"amount"`
}
