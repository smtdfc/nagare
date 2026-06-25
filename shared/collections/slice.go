package nagare_collections

func FindStringIntersection(slice1 []string, slice2 []string) []string {
	// Create a set from the first slice
	elements := make(map[string]bool)
	for _, str := range slice1 {
		elements[str] = true
	}

	intersection := []string{}
	// Check against the second slice
	for _, str := range slice2 {
		if elements[str] {
			intersection = append(intersection, str)
			// Remove to avoid adding the same string again
			delete(elements, str)
		}
	}

	return intersection
}

func HasStringIntersection(slice1 []string, slice2 []string) bool {
	// Create a set from the first slice
	elements := make(map[string]bool)
	for _, str := range slice1 {
		elements[str] = true
	}

	// Check against the second slice
	for _, str := range slice2 {
		if elements[str] {
			// Return true immediately as soon as a common element is found
			return true
		}
	}

	return false
}
