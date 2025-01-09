package api

// urlConverter defines the rules of the converting the data for URL-correct format
type URLConverter struct{}

// getFilters converts filters in the map format.
func (c URLConverter) GetFilters(filters []string) map[string]string {
	var filtersURL = make(map[string]string)

	for i := 1; i < len(filters); i += 2 {
		filtersURL[filters[i-1]] = filters[i]
	}

	return filtersURL
}
