package dto

type SendCoinRequest struct {
	ToUser string
	Amount int
}

type TransactionDetail struct {
	FromUser string
	ToUser   string
	Amount   int
}

type CoinHistoryResponse struct {
	Received []TransactionDetail
	Sent     []TransactionDetail
}
