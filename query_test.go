package wikipedia

import (
	"log"
	"testing"
)

func TestListAllPages(t *testing.T) {
	response, err := ListAllPages(ListAllPagesRequest{
		From:  "A",
		To:    "B",
		Limit: 20,
	})
	if err != nil {
		t.Fail()
		return
	}
	log.Print(response)
}
