package wikipedia

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// RedirectType to be utilised as a redirect filter in the `allpages` query
type RedirectType string

const All RedirectType = "all"
const Redirects RedirectType = "redirects"
const NonRedirects RedirectType = "nonredirects"

// ProtectionType to be utilised as a protection type filter in the `allpages` query
type ProtectionType string

const Edit ProtectionType = "edit"
const Move ProtectionType = "move"
const Upload ProtectionType = "upload"

type ListAllPagesRequest struct {
	//The page title to start enumerating from.
	From string `json:"apfrom"`
	//When more results are available, use this to continue.
	Continue string `json:"apcontinue"`
	//The page title to stop enumerating at.
	To string `json:"apto"`
	//Search for all page titles that begin with this value.
	Prefix string `json:"apprefix"`
	//The namespace to enumerate.
	//One of the following values: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 90, 91, 92, 93, 100, 101, 102, 103, 104, 105, 106, 107, 486, 487, 828, 829, 1198, 1199, 2300, 2301, 2302, 2303, 2600, 5500, 5501
	//Default: 0
	Namespace int64 `json:"apnamespace"`
	//Which pages to list.
	//Note: Due to miser mode, using this may result in fewer than aplimit results returned before continuing; in extreme cases, zero results may be returned.
	//One of the following values: all, redirects, nonredirects
	//Default: all
	FilterRedir RedirectType `json:"apfilterredir"`
	//Limit to pages with at least this many bytes.
	//Type: integer
	MaxSize string `json:"apmaxsize"`
	//
	//Limit to pages with at most this many bytes.
	//Type: integer
	MinSize string `json:"apminsize"`
	//Limit to protected pages only.
	//Values (separate with | or alternative): edit, move, upload
	ProtectionType ProtectionType `json:"apprtype"`
	//Filter protections based on protection level (must be used with apprtype= parameter).
	//Values (separate with | or alternative): Can be empty, or autoconfirmed, sysop
	PrLevel string `json:"apprlevel"`
	//Filter protections based on cascadingness (ignored when apprtype isn't set).
	//One of the following values: cascading, noncascading, all
	//Default: all
	PrFilterCascade string `json:"apprfiltercascade"`
	//How many total pages to return.
	//No more than 500 (5,000 for bots) allowed.
	//Type: integer or max
	//Default: 10
	Limit int64 `json:"aplimit"`
	//The direction in which to list.
	//One of the following values: ascending, descending
	//Default: ascending
	Direction Direction `json:"apdir"`
	//Filter based on whether a page has langlinks. Note that this may not consider langlinks added by extensions.
	//One of the following values: withlanglinks, withoutlanglinks, all
	//Default: all
	FilterLangLinks string `json:"apfilterlanglinks"`
	//Which protection expiry to filter the page on:
	//indefinite
	//Get only pages with indefinite protection expiry.
	//definite
	//Get only pages with a definite (specific) protection expiry.
	//all
	//Get pages with any protections expiry.
	//One of the following values: indefinite, definite, all
	//Default: all
	PrExpiry string `json:"apprexpiry"`
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
	query.Add("aplimit", strconv.FormatInt(request.Limit, 10))
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
