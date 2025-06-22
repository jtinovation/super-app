package dto

type QueryParams struct {
	Page    int
	PerPage int
	Search  string
	Sort    string
	Order   string
	Filter  map[string]string
}