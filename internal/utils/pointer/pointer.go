// Package pointer generates pointers to values in simple one liners without requiring the declaration of a
// separate variable.
package pointer

// Pointer

// ToInt generates a pointer to an int.
func ToInt(i int) *int {
	return &i
}
