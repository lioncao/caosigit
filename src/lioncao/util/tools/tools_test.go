package tools

import (
	"math"
	"testing"
)

func Test_main(t *testing.T) {
	// _test_math()
}

func _test_math() {

	count := 100000
	list := make([]int64, count)

	for i := 0; i < count; i++ {
		list[i] = int64(i)
	}

	count = 1
	for i := 0; i < count; i++ {
		randomValue := int64(49999) // int64(GetRand2(0, 1000))
		idx := Math_valueToIndex(list, randomValue)
		ShowInfo(i, randomValue, idx)

	}

}

// func FastInvSqrt(float64 x) float64 {
//    xhalf := float64(0.5f * x)
//    i := int64(x)    // evil floating point bit level hacking
//   i = 0x5f3759df - (i >> 1);  // what the fuck?
//   x = *(float*)&i;
//   x = x*(1.5f-(xhalf*x*x));
//   return x
// }

func _test_sqrt() {
	// math.Sqrt2
}

// float FastInvSqrt(float x) {
//   float xhalf = 0.5f * x;
//   int i = *(int*)&x;         // evil floating point bit level hacking
//   i = 0x5f3759df - (i >> 1);  // what the fuck?
//   x = *(float*)&i;
//   x = x*(1.5f-(xhalf*x*x));
//   return x;
// }
