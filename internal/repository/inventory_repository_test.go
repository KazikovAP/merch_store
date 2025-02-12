package repository_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/repository"
)

func TestGetByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "item_type", "quantity"}).
		AddRow(1, 1, "t-shirt", 2).
		AddRow(2, 1, "pen", 5)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, item_type, quantity FROM inventory WHERE user_id=$1")).
		WithArgs(1).
		WillReturnRows(rows)

	invRepo := repository.NewInventoryRepository(db)

	items, err := invRepo.GetByUserID(1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 rows, got %d", len(items))
	}
}

func TestAddItem_InsertNew(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM inventory WHERE user_id=$1 AND item_type=$2")).
		WithArgs(1, "t-shirt").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO inventory (user_id, item_type, quantity) VALUES ($1, $2, $3)")).
		WithArgs(1, "t-shirt", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	invRepo := repository.NewInventoryRepository(db)

	err = invRepo.AddItem(1, "t-shirt", 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestAddItem_UpdateExisting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM inventory WHERE user_id=$1 AND item_type=$2")).
		WithArgs(1, "t-shirt").
		WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE inventory SET quantity = quantity + $1 WHERE id=$2")).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	invRepo := repository.NewInventoryRepository(db)

	err = invRepo.AddItem(1, "t-shirt", 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
