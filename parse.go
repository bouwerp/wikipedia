package wikipedia

type ParseParams struct {
	Params
	Page   string `json:"page"`
	Format string `json:"format"`
}
