package wikipedia

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ListAllPagesRequest struct {
	From  string
	To    string
	Limit int
	// this field, if populated
	Continue string
}

type Page struct {
	Pageid int    `json:"pageid"`
	Ns     int    `json:"ns"`
	Title  string `json:"title"`
}

type ListAllPagesResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		// The value of apcontinue must be used in the next request's
		// Continue field, until Continue is empty (will have '-||' when there is still pages)
		Apcontinue string `json:"apcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Allpages []Page `json:"allpages"`
	} `json:"query"`
}

func (r ListAllPagesRequest) validate() error {
	if r.Limit > 500 {
		return LimitTooHigh{}
	}
	return nil
}

func ListAllPages(request ListAllPagesRequest) (*ListAllPagesResponse, error) {
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
	query.Add("list", "allpages")
	query.Add("format", "json")
	query.Add("apfrom", request.From)
	query.Add("apto", request.To)
	query.Add("aplimit", strconv.Itoa(request.Limit))
	query.Add("apcontinue", request.Continue)
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
	var response ListAllPagesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &response, nil
}
