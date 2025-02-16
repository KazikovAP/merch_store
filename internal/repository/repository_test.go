package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KazikovAP/merch_store/internal/repository"
	"github.com/stretchr/testify/assert"
)

// Тест на успешную инициализацию репозиториев.
func TestNewRepositories_Success(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repos := repository.NewRepositories(db)

	assert.NotNil(t, repos.User)
	assert.NotNil(t, repos.Transaction)
	assert.NotNil(t, repos.Purchase)
	assert.NotNil(t, repos.Merch)
}
