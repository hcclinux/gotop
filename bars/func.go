package bars

import (
	"sort"
)


// BarsCount 求总周期数
func BarsCount(dataset []float64) int {
	return len(dataset)
}


// Ref 引用period周期前的数据
func Ref(dataset []float64, period int) float64 {
	if period < len(dataset) {
		return dataset[len(dataset) - period - 1]
	}
	return dataset[len(dataset) - 1]
}

// TopRange 当前值是近多少周期内的最大值
func TopRange(dataset []float64, val float64) int {
	period := 0
	if len(dataset) < 1 {
		return period
	}

	for i := len(dataset) - 2; i >= 0; i-- {
		if dataset[i] >= val {
			return period
		}
		period++
	}
	return period
}


// LowRange 当前值是近多少周期内的最小值
func LowRange(dataset []float64, val float64) int {
	period := 0
	if len(dataset) < 1 {
		return period
	}

	for i := len(dataset) - 2; i >= 0; i-- {
		if dataset[i] <= val {
			return period
		}
		period++
	}
	return period
}

// HHV 求最大值
func HHV(dataset []float64, period int) float64 {
	var length int
	if period == 0 {
		length = len(dataset)
	} else {
		length = period
	}

	temp := make([]float64, length)
	switch period {
	case 0:
		copy(temp, dataset)
	default:
		copy(temp, dataset[len(dataset) - period:])
	}

	sort.Float64s(temp)
	return temp[len(temp) - 1]
}

// LLV 求最小值
func LLV(dataset []float64, period int) float64 {
	var length int
	if period == 0 {
		length = len(dataset)
	} else {
		length = period
	}
	temp := make([]float64, length)
	switch period {
	case 0:
		copy(temp, dataset)
	default:
		copy(temp, dataset[len(dataset) - period:])
	}

	sort.Float64s(temp)
	return temp[0]
}

// HHVBars 求上一高点到当前周期数
func HHVBars(dataset []float64, period int) int {
	if period == 1 {
		return 0
	}
	if period != 0 {
		dataset = dataset[len(dataset)-period:]
	}
	val := HHV(dataset, period)

	idxList := DuplicateElement(dataset, val)
	index := idxList[len(idxList)-1]
	return len(dataset[index+1:])
}


// LLVBars 求上一低点到当前周期数
func LLVBars(dataset []float64, period int) int {
	if period == 1 {
		return 0
	} else if period != 0 {
		dataset = dataset[len(dataset)-period:]
	}
	val := LLV(dataset, period)
	idxList := DuplicateElement(dataset, val)
	index := idxList[len(idxList)-1]
	return len(dataset[index+1:])
}


// DuplicateElement 查找重复元素下标
func DuplicateElement(dataset []float64, val float64) []int {
	temp := make(map[float64]struct{})
	temp[val] = struct{}{}
	indexList := make([]int, 0)
	for i := 0; i < len(dataset); i++ {
        if _, ok := temp[dataset[i]]; ok {
            indexList = append(indexList, i)
		}
	}
	return indexList
}

// Count 统计满足条件的周期数
func Count(dataset []float64, period int, f func(x, y float64)(bool)) int {
	count := 0
	l := len(dataset)
	for i := l - 1; i > l - period - 1; i-- {
		if f(dataset[i], dataset[i-1]) {
			count++
		}
	}
	return count
}

// CountC 统计
func CountC(low, high []float64, period int, f ...func(x, y float64)(bool)) int {
	count := 0
	l := len(low)
	for i := l - 1; i > l - period - 1; i-- {
		dbx := f[1](low[i], low[i-1]) && f[0](high[i], high[i-1])
		xbd := f[0](low[i], low[i-1]) && f[1](high[i], high[i-1])
		if dbx || xbd {
			count++
		}
	}
	return count
}


// Min .
func Min(val1, val2 float64) float64 {
	if val1 < val2 {
		return val1
	}
	return val2
}

// Max .
func Max(val1, val2 float64) float64 {
	if val1 > val2 {
		return val1
	}
	return val2
}