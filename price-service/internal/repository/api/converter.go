package api

// URLConverter defines the rules of the converting the data for URL-correct format
type urlConverter struct{}

// getFilters converts filters in the map format.
func (url urlConverter) getFilters(filters []string) map[string]string {
	var filtersURL = make(map[string]string)

	for i := 1; i < len(filters); i += 2 {
		filtersURL[filters[i-1]] = filters[i]
	}

	return filtersURL
}
