package wikipedia

const ApiUrl = "https://www.mediawiki.org/w/api.php"

type Params struct {
	Action string `json:"action"`
}
