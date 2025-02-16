package merch_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	"github.com/KazikovAP/merch_store/internal/repository/merch"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное получение товара по имени.
func TestRepo_GetByName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	name := "t-shirt"
	expectedMerch := &domain.Merchandise{Name: name, Price: 80}
	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"name", "price"}).
			AddRow(expectedMerch.Name, expectedMerch.Price))

	repo := merch.NewMerchRepository(db)

	m, err := repo.GetByName(context.Background(), name)
	assert.NoError(t, err)
	assert.Equal(t, expectedMerch, m)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении товара по имени (товар не найден).
func TestRepo_GetByName_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	name := "unknown-item"
	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)

	repo := merch.NewMerchRepository(db)

	m, err := repo.GetByName(context.Background(), name)
	assert.Error(t, err)
	assert.Nil(t, m)
	assert.Contains(t, err.Error(), "failed to get merchandise by name")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на успешное получение списка товаров.
func TestRepo_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	expectedItems := []domain.Merchandise{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
	}

	rows := sqlmock.NewRows([]string{"name", "price"})
	for _, item := range expectedItems {
		rows.AddRow(item.Name, item.Price)
	}

	mock.ExpectQuery(`SELECT name, price FROM merchandise`).
		WillReturnRows(rows)

	repo := merch.NewMerchRepository(db)

	items, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при получении списка товаров.
func TestRepo_List_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT name, price FROM merchandise`).
		WillReturnError(errors.New("database error"))

	repo := merch.NewMerchRepository(db)

	items, err := repo.List(context.Background())
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Contains(t, err.Error(), "failed to query merchandise")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
