package conversions

/* conversions: the package contains pure functions that convert the values of the input
a different type of the same structure, eg. a list.
*/

import "strconv"

// StringsToInt converts the list of strings into list of integers.
func StringsToInt(x []string) ([]int, error) {
	r := make([]int, len(x))
	for c, i := range x {
		j, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		r[c] = j
	}
	return r, nil
}

// BytesToMb returns the mb size of the input number
func BytesToMb(b float64) float64 {
	return b / 1048576
}
