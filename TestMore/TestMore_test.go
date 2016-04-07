package TestMore

import "testing"

import "gmon/TestMore"

func TestOK(t *testing.T) {
	t.Error(TestMore.Ok(1, 1, 1, 1, "ALL ONE") == true)
}
