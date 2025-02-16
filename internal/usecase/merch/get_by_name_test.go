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

// Тест на успешное получение товара по имени.
func TestUseCase_GetByName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	name := "t-shirt"
	expectedMerch := &domain.Merchandise{
		Name:  name,
		Price: 80,
	}

	rows := sqlmock.NewRows([]string{"name", "price"}).
		AddRow(expectedMerch.Name, expectedMerch.Price)
	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(name).
		WillReturnRows(rows)

	repo := merRep.NewMerchRepository(db)
	useCase := merch.NewUseCase(repo)

	mrch, err := useCase.GetByName(context.Background(), name)
	assert.NoError(t, err)
	assert.Equal(t, expectedMerch, mrch)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// Тест на ошибку при пустом имени товара.
func TestUseCase_GetByName_EmptyName(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := merRep.NewMerchRepository(db)
	useCase := merch.NewUseCase(repo)

	mrch, err := useCase.GetByName(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, mrch)
	assert.Contains(t, err.Error(), "empty merchandise name")
}

// Тест на ошибку репозитория.
func TestUseCase_GetByName_RepoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	name := "t-shirt"
	repoErr := errors.New("database error")

	mock.ExpectQuery(`SELECT name, price FROM merchandise WHERE name = \$1`).
		WithArgs(name).
		WillReturnError(repoErr)

	repo := merRep.NewMerchRepository(db)
	useCase := merch.NewUseCase(repo)

	mrch, err := useCase.GetByName(context.Background(), name)
	assert.Error(t, err)
	assert.Nil(t, mrch)
	assert.Contains(t, err.Error(), "failed to get merchandise by name")
}
