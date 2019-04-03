package wikipedia

import (
	"log"
	"strconv"
	"testing"
)

func TestListAllPages(t *testing.T) {
	// general pattern for getting all pages from A to B in batches of 200
	var pages []Page
	apContinue := ""
	for {
		response, err := ListAllPages(ListAllPagesRequest{
			From:     "A",
			To:       "B",
			Limit:    200,
			Continue: apContinue,
		})
		if err != nil {
			t.Fail()
			return
		}
		pages = append(pages, response.Query.Allpages...)
		if response.Continue.Continue == "" {
			break
		}
		apContinue = response.Continue.Apcontinue
	}
	log.Println("Total: Found " + strconv.Itoa(len(pages)) + " pages")
}
