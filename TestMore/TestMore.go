package TestMore

import "fmt"

// Ok takes a pair of values, a pair of expected values and a message string, return bool of success.
func Ok(ret interface{}, err interface{}, exp interface{}, expErr interface{}, message string) bool {
	ok := false
	errOk := true
	if ret == exp {
		message = message + " return OK"
	}
	if err == expErr {
		message = message + " error OK"
	}

	fmt.Println(message)

	return ok && errOk

}
