package util

// Check if a specific string (needle) exist in a specific slice (haystack)
func Contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}
	}
	return false
}
