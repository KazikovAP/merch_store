package merch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/domain"
	merRep "github.com/KazikovAP/merch_store/internal/repository/merch"
	"github.com/KazikovAP/merch_store/internal/usecase/merch"
	"github.com/stretchr/testify/assert"
)

// Тест на успешное получение списка товаров.
func TestUseCase_List_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	expectedMerch := []domain.Merchandise{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
	}

	rows := sqlmock.NewRows([]string{"name", "price"})
	for _, item := range expectedMerch {
		rows.AddRow(item.Name, item.Price)
	}

	mock.ExpectQuery(`SELECT name, price FROM merchandise`).
		WillReturnRows(rows)

	repo := merRep.NewMerchRepository(db)
	useCase := merch.NewUseCase(repo)

	mrch, err := useCase.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedMerch, mrch)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку репозитория.
func TestUseCase_List_RepoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repoErr := errors.New("database error")

	mock.ExpectQuery(`SELECT name, price FROM merchandise`).
		WillReturnError(repoErr)

	repo := merRep.NewMerchRepository(db)
	useCase := merch.NewUseCase(repo)

	mrch, err := useCase.List(context.Background())
	assert.Error(t, err)
	assert.Nil(t, mrch)
	assert.Contains(t, err.Error(), "failed to get merchandise list")
}
