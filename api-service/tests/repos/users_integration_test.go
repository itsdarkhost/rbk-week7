package repos_test

import (
	"context"
	"testing"

	"github.com/itsdarkhost/rbk-week4/internal/repos"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestUserRepoIntegrationCreateAndRead(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			deleted_at TIMESTAMP NULL
		);
	`)
	require.NoError(t, err)

	repo := repos.NewUserRepo(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, "Alice", "alice@example.com", "hash", "admin")
	require.NoError(t, err)
	require.NotZero(t, created.Id)

	byID, err := repo.Get(ctx, created.Id)
	require.NoError(t, err)
	assert.Equal(t, created.Email, byID.Email)
	assert.Equal(t, "admin", byID.Role)

	byEmail, err := repo.GetByEmail(ctx, "alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, created.Id, byEmail.Id)

	users, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, "Alice", users[0].Username)
}
