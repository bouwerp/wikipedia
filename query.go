package wikipedia

import (
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
}

type ListAllPagesResponse struct {
}

func ListAllPages(request ListAllPagesRequest) (*ListAllPagesResponse, error) {
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
	u.RawQuery = query.Encode()

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

	log.Println(string(responseBytes))

	return &ListAllPagesResponse{}, nil
}
