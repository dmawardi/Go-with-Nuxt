package utility

// Function to sort permissions data from enforcer
func SortMapStringInterfaceAlphabetically(a, b map[string]interface{}, key string) bool {
	resourceA, okA := a[key].(string)
	resourceB, okB := b[key].(string)

	// If either of the elements doesn't have a valid "resource" string, consider it greater (move it to the end)
	if !okA || !okB {
		return false
	}

	// Compare the "resource" strings alphabetically
	return resourceA < resourceB
}

// Function to sort a map[string]string by a given key
func SortMapStringStringAlphabetically(a, b map[string]string, key string) bool {
	valueA, okA := a[key]
	valueB, okB := b[key]

	// If either of the elements doesn't have a valid string for the given key, consider it greater (move it to the end)
	if !okA || !okB {
		return false
	}

	// Compare the strings alphabetically
	return valueA < valueB
}
