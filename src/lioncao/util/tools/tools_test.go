package tools

import (
	"testing"
)

func Test_main(t *testing.T) {
	// _test_math()
}

func _test_math() {
	list := []int64{
		4, 5, 6, 7, 8,
		9, 10, 11, 12, 13,
		14, 15, 16, 17, 18,
		19, 20, 21, 22, 23,
		24, 25, 28, 34, 35,
		36, 38, 43, 56, 87,
	}

	value := int64(16)
	idx := Math_valueToIndex(list, value)

	ShowInfo("...", value, idx)

}
