package sqlite

import (
	"arshin/internal/domain"
	"context"
	"database/sql"
	"errors"
	"log"
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
	var errdata domain.ErrInvalidData
	if errors.As(err, &errdata) {
		t.Fatal(res, err)
	}
	log.Println(res)
}
