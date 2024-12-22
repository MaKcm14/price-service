package api

// URLConverter defines the rules of the converting the data for URL-correct format
type urlConverter struct{}

// convertFilters converts the filters to the url view (divided by '=': &key1=value1&key2=value2&...).
func (url urlConverter) convertFilters(filters []string) string {
	var filtersURL string

	for i := 1; i < len(filters); i += 2 {
		filtersURL += "&" + filters[i-1] + "=" + filters[i]
	}

	return filtersURL
}

// getFilters converts filters in the map format.
func (url urlConverter) getFilters(filters []string) map[string]string {
	var filtersURL = make(map[string]string)

	for i := 1; i < len(filters); i += 2 {
		filtersURL[filters[i-1]] = filters[i]
	}

	return filtersURL
}
