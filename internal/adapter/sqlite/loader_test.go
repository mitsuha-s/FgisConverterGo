package sqlite

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/gwenn/gosqlite"
)

func TestLoader(t *testing.T) {
	db, err := sql.Open("sqlite3", "/home/mitsuha/test.db")
	if err != nil {
		t.Fatal(err)
	}
	repo := New(db)
	res, err := repo.Load(context.Background())
	if err != nil {
		t.Fatal(res, err)
	}
}
