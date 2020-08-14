// +build integration

package user

import (
	"context"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/require"
)

func setUp(t *testing.T) *sqlx.DB {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		t.Fatal("DATABASE_DSN is empty")
	}

	db := sqlx.MustConnect("pgx", dsn)

	err := goose.Up(db.DB, "../../migrations")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func tearDown(t *testing.T, db *sqlx.DB) {
	db.MustExec("DELETE FROM public.user WHERE id != 1")
	db.Close()
}

func TestManager_GetUserByEmailPassword(t *testing.T) {
	db := setUp(t)
	defer tearDown(t, db)
	m := NewManager(db)

	u, err := m.GetUserByEmailPassword(context.Background(), "admin", "admin")
	require.NoError(t, err)
	require.NotEmpty(t, u)
}

func TestManager_Blacklist(t *testing.T) {
	db := setUp(t)
	defer tearDown(t, db)
	m := NewManager(db)

	ctx := context.Background()
	err := m.BlacklistUserToken(ctx, 1, "token")
	require.NoError(t, err)

	blacklisted := m.IsTokenBlacklisted(ctx, 1, "token")
	require.True(t, blacklisted)

	blacklisted = m.IsTokenBlacklisted(ctx, 1, "token2")
	require.False(t, blacklisted)
}

func TestManager_Create(t *testing.T) {
	db := setUp(t)
	defer tearDown(t, db)
	m := NewManager(db)

	var testCases = []struct {
		email    string
		password string
		phone    string
	}{
		{"email@domain.com", "random", "phone1"},
		{"verylooooongemail@veryyyyloooooongdomain.com", "random", "phon21"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			ctx := context.Background()
			createdUser, err := m.Create(ctx, "name", tc.email, tc.password, tc.phone, RoleUser)
			require.NoError(t, err)

			u, err := m.GetUserByEmailPassword(ctx, tc.email, tc.password)
			require.NoError(t, err)

			require.Equal(t, createdUser.ID, u.ID)
		})
	}

}
