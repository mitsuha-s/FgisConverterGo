package xlsx

import (
	"log"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	parser := NewParser(log.New(os.Stdout, "TEST ", 0))
	_, err := parser.Load("/home/mitsuha/муром октябрь 05.11.xlsx")
	log.Println(err)
}
