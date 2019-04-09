package wikipedia

const ApiUrl = "https://www.mediawiki.org/w/api.php"

type Params struct {
	Action string `json:"action"`
}

// Direction sort direction
type Direction string

// ASC ascending sort direction
const ASC Direction = "ascending"

// DESC descending sort direction
const DESC Direction = "descending"

// Property
type Property string

const SIZE Property = "size"
const HIDDEN Property = "hidden"
