package xml

import (
	"arshin/internal/domain"
	"testing"
)

func TestConverter(t *testing.T) {
	v := domain.NewTestVirifications()

	c := converter{}
	err := c.Deliver("./test.out", v)
	if err != nil {
		t.Fatal(err)
	}
}
