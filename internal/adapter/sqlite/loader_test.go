package sqlite

import (
	"arshin/internal/domain"
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	_ "github.com/gwenn/gosqlite"
)

func TestLoader(t *testing.T) {
	db, err := sql.Open("sqlite3", "/home/mitsuha/test.db")
	if err != nil {
		t.Fatal(err)
	}
	repo := New(db)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	res, err := repo.GetConfig(ctx)
	var errdata domain.ErrInvalidData
	if errors.As(err, &errdata) {
		t.Fatal(res, err)
	}
	log.Println(res, err)
}
