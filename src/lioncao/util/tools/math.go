package tools

// 根据数值查找index
func Math_valueToIndex(list []int64, value int64) (idx int) {

	var (
		half  int
		vhalf int64
	)

	curList := list
	cnt := len(curList)

	if cnt <= 1 {
		return 0
	}

	begin := 0
	for {
		ShowInfo("===", value, begin, cnt)
		if cnt <= 6 {
			for i := 1; i < cnt; i++ {
				if value < curList[i] {
					return begin + i - 1
				}
			}
			return begin + cnt - 1
		}

		half = cnt / 2
		vhalf = curList[half]

		if value == vhalf {
			return begin + half
		} else if value < vhalf {
			curList = curList[:half]
			cnt = half
		} else {
			curList = curList[half:]
			cnt = cnt - half
			begin = begin + half
		}

	}

	return idx
}
