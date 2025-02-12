package service_test

import (
	"errors"

	"github.com/KazikovAP/merch_store/internal/model"
)

type mockUserRepo struct {
	user *model.User
}

func (m *mockUserRepo) GetByUsername(username string) (*model.User, error) {
	if m.user != nil && m.user.Username == username {
		return m.user, nil
	}

	return nil, errors.New("user not found")
}

func (m *mockUserRepo) Create(user *model.User) error {
	m.user = user
	return nil
}

func (m *mockUserRepo) Update(user *model.User) error {
	m.user = user
	return nil
}

type mockUserRepoMap struct {
	users map[string]*model.User
}

func (m *mockUserRepoMap) GetByUsername(username string) (*model.User, error) {
	user, ok := m.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (m *mockUserRepoMap) Create(user *model.User) error {
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepoMap) Update(user *model.User) error {
	m.users[user.Username] = user
	return nil
}

type mockTransactionRepo struct {
	transactions []*model.Transaction
}

func (m *mockTransactionRepo) Create(txn *model.Transaction) error {
	m.transactions = append(m.transactions, txn)
	return nil
}

func (m *mockTransactionRepo) GetByUserID(userID int) ([]*model.Transaction, error) {
	var res []*model.Transaction

	for _, txn := range m.transactions {
		if txn.UserID == userID {
			res = append(res, txn)
		}
	}

	return res, nil
}

type mockInventoryRepo struct {
	items []*model.InventoryItem
}

func (m *mockInventoryRepo) GetByUserID(userID int) ([]*model.InventoryItem, error) {
	var res []*model.InventoryItem

	for _, item := range m.items {
		if item.UserID == userID {
			res = append(res, item)
		}
	}

	return res, nil
}

func (m *mockInventoryRepo) AddItem(userID int, itemType string, quantity int) error {
	for _, item := range m.items {
		if item.UserID == userID && item.ItemType == itemType {
			item.Quantity += quantity
			return nil
		}
	}

	m.items = append(m.items, &model.InventoryItem{
		ID:       len(m.items) + 1,
		UserID:   userID,
		ItemType: itemType,
		Quantity: quantity,
	})

	return nil
}
