package repository_test

import (
	"library/internal/db"
	"library/internal/models"
	"library/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndSelectByIdAuthor_When_OK(t *testing.T) {
	db, err := db.InitDb(":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NotEmpty(t, db)
	store := repository.NewLibraryStore(db)

	author := models.Author{
		Name:    "Пупс Пупсов",
		Country: "Russia",
	}

	id, err := store.CreateAuthor(author)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	author.Id = id

	actualAuthor, err := store.GetAuthorById(id)
	require.NoError(t, err)
	require.NotEmpty(t, actualAuthor)

	assert.Equal(t, author, actualAuthor)
}
