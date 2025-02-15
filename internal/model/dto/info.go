package dto

type InfoResponse struct {
	Coins       int
	Inventory   []InventoryResponse
	CoinHistory CoinHistoryResponse
}
