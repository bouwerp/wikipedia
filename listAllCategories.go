package wikipedia

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ListAllCategoriesRequest struct {
	// The category to start enumerating from.
	From string `json:"acfrom"`
	// When more results are available, use this to continue.
	Continue string `json:"accontinue"`
	//The category to stop enumerating at.
	To string `json:"acto"`
	//Search for all category titles that begin with this value.
	Prefix string `json:"acprefix"`
	// Direction to sort in.
	// One of the following values: ascending, descending
	// Default: ascending
	Dir Direction `json:"acdir"`
	// Only return categories with at least this many members.
	// Type: integer
	Min int64 `json:"acmin"`
	// Only return categories with at most this many members.
	// Type: integer
	Max int64 `json:"acmax"`
	// How many categories to return.
	// No more than 500 (5,000 for bots) allowed.
	// Type: integer or max
	// Default: 10
	Limit int64 `json:"aclimit"`
	//Which properties to get:
	//size
	//Adds number of pages in the category.
	//hidden
	//Tags categories that are hidden with __HIDDENCAT__.
	//Values (separate with | or alternative): size, hidden
	//Default: (empty)
	Prop []Property `json:"acprop"`
}

type ListAllCategoriesResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Accontinue string `json:"accontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Categories []Category `json:"allcategories"`
	} `json:"query"`
}

type Category struct {
	Size    int    `json:"size"`
	Pages   int    `json:"pages"`
	Files   int    `json:"files"`
	Subcats int    `json:"subcats"`
	Name    string `json:"*"`
}

func (r ListAllCategoriesRequest) validate() error {
	if r.Limit > 500 {
		return LimitTooHigh{}
	}
	return nil
}

func ListAllCategories(request ListAllCategoriesRequest) (*ListAllCategoriesResponse, error) {
	// validate request
	if err := request.validate(); err != nil {
		return nil, err
	}

	// construct URL
	u, err := url.Parse(ApiUrl)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Add("action", "query")
	query.Add("list", "allcategories")
	query.Add("format", "json")
	query.Add("acfrom", request.From)
	query.Add("accontinue", request.Continue)
	query.Add("acto", request.To)
	query.Add("acmin", strconv.FormatInt(request.Min, 10))
	query.Add("acprefix", request.Prefix)
	query.Add("acdir", string(request.Dir))
	query.Add("acmax", strconv.FormatInt(request.Max, 10))
	query.Add("aclimit", strconv.FormatInt(request.Limit, 10))
	var props []string
	for _, p := range request.Prop {
		props = append(props, string(p))
	}
	query.Add("acprop", strings.Join(props, "|"))

	u.RawQuery = query.Encode()

	// execute the request
	resp, err := http.Get(u.String())
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	var response ListAllCategoriesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &response, nil
}
