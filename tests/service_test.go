package repository_test

import (
	"library/internal/db"
	"library/internal/models"
	"library/internal/repository"
	"library/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceCreateAuthor_When_InvalidName(t *testing.T) {
	db, err := db.InitDb(":memory:")
	require.NoError(t, err)
	defer db.Close()

	store := repository.NewLibraryStore(db)
	svc := service.NewLibraryService(store)

	_, err = svc.CreateAuthor(models.Author{
		Name:    "",
		Country: "Russia",
	})
	require.Error(t, err)
}
