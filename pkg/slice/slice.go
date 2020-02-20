// contains some helper methods for working with slices.
package slice

// function to check if given slice contains
// given object.
// It does not do a deep check, but simply comparing them
// with "=="
//
// only checks type "string"
func ContainsString(sl []string, obj string) bool {
	for _, o := range sl {
		if o == obj {
			return true
		}
	}
	return false
}

// function to check if given slice contains
// given object.
// It does not do a deep check, but simply comparing them
// with "=="
//
// only checks type "int"
func ContainsInt(sl []int, obj int) bool {
	for _, o := range sl {
		if o == obj {
			return true
		}
	}
	return false
}
