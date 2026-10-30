package xls

import (
	"log"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	parser := New(log.New(os.Stdout, "TEST ", 0))
	res, err := parser.Load("/home/mitsuha/муром октябрь 05.11.xlsx")
	log.Println(res, err)
}
