package cvmath

import (
	"math"
	"sort"

	"github.com/digisan/go-generics/i64"
	"github.com/digisan/go-generics/u8i64"
)

func Histogram(data ...byte) (m map[byte]int, maxIdx byte, maxN int) {
	m = make(map[byte]int)
	for i := 0; i < 256; i++ {
		m[byte(i)] = 0
	}
	for i := 0; i < len(data); i++ {
		v := data[i]
		m[v]++
		if m[v] > maxN {
			maxIdx, maxN = v, m[v]
		}
	}
	return
}

func fnIsPeak(data []float64, halfstep int) func(i int) bool {

	var (
		deri  = Derivative1(data...)
		E     = 0.1
		iPeak = []int{}
	)

	return func(i int) bool {
		if i >= halfstep && i < len(deri)-halfstep {

			nUp1, nEven1, nDown1 := 0, 0, 0
			for j := i - halfstep; j < i; j++ {
				switch {
				case deri[j] > E: // up
					nUp1++
				case deri[j] >= -E && deri[j] <= E: // top
					nEven1++
				default: // down
					nDown1++
				}
			}
			if nDown1 > 0 || (nUp1 <= nEven1) {
				return false
			}

			nUp2, nEven2, nDown2 := 0, 0, 0
			for j := i + halfstep; j > i; j-- {
				switch {
				case deri[j] < -E: // down
					nDown2++
				case deri[j] >= -E && deri[j] <= E: // top
					nEven2++
				default: // up
					nUp2++
				}
			}
			if nUp2 > 0 || (nDown2 <= nEven2) {
				return false
			}

			_, e := MaxMinAbs(deri[i-halfstep : i+halfstep+1]...)
			if math.Abs(deri[i]) == e {
				if len(iPeak) > 0 && i <= iPeak[len(iPeak)-1]+halfstep {
					return false
				}
				iPeak = append(iPeak, i)
				return true
			}
		}
		return false
	}
}

func fnIsBottom(data []float64, halfstep int) func(i int) bool {

	var (
		deri    = Derivative1(data...)
		E       = 0.1
		iBottom = []int{}
	)

	return func(i int) bool {
		if i >= halfstep && i < len(deri)-halfstep {

			nUp1, nEven1, nDown1 := 0, 0, 0
			for j := i - halfstep; j < i; j++ {
				switch {
				case deri[j] < -E: // down
					nDown1++
				case deri[j] >= -E && deri[j] <= E: // bottom
					nEven1++
				default: // up
					nUp1++
				}
			}
			if nUp1 > 0 || (nDown1 <= nEven1) {
				return false
			}

			nUp2, nEven2, nDown2 := 0, 0, 0
			for j := i + halfstep; j > i; j-- {
				switch {
				case deri[j] > E: // up
					nUp2++
				case deri[j] >= -E && deri[j] <= E: // bottom
					nEven2++
				default: // down
					nDown2++
				}
			}
			if nDown2 > 0 || (nUp2 <= nEven2) {
				return false
			}

			_, e := MaxMinAbs(deri[i-halfstep : i+halfstep+1]...)
			if math.Abs(deri[i]) == e {
				if len(iBottom) > 0 && i <= iBottom[len(iBottom)-1]+halfstep {
					return false
				}
				iBottom = append(iBottom, i)
				return true
			}
		}
		return false
	}
}

func Peaks(data map[byte]int, halfstep, nSmooth, nPeak int) map[byte]int {

	m := make(map[byte]int)
	ks, vs := u8i64.Map2KVs(data, func(i, j byte) bool { return i < j }, nil)
	vsTemp := I32ToF64(vs)
	for i := 0; i < nSmooth; i++ {
		vsTemp = smooth9(vsTemp)
	}

	isPeak := fnIsPeak(vsTemp, halfstep)
	for i := 0; i < len(vs); i++ {
		if isPeak(i) {
			m[ks[i]] = vs[i]
		}
	}

	// adjust to max value
	mp := make(map[byte]int)
	for k, v := range m {
		if max, n := i64.MaxIdx(vs[k-1 : k+2]...); max > v {
			mp[k-1+byte(n)] = max
		} else {
			mp[k] = v
		}
	}

	if nPeak < 0 {
		return mp
	}

	ks, vs = u8i64.Map2KVs(mp, nil, func(i, j int) bool { return j < i })
	mp = make(map[byte]int)
	for i := 0; i < nPeak && i < len(ks); i++ {
		mp[ks[i]] = vs[i]
	}
	return mp
}

func Bottoms(data map[byte]int, halfstep, nSmooth, nBottom int) map[byte]int {

	m := make(map[byte]int)
	ks, vs := u8i64.Map2KVs(data, func(i, j byte) bool { return i < j }, nil)
	vsTemp := I32ToF64(vs)
	for i := 0; i < nSmooth; i++ {
		vsTemp = smooth9(vsTemp)
	}

	isBottom := fnIsBottom(vsTemp, halfstep)
	for i := 0; i < len(vs); i++ {
		if isBottom(i) {
			m[ks[i]] = vs[i]
		}
	}

	// adjust to min value
	mp := make(map[byte]int)
	for k, v := range m {
		if min, n := i64.MinIdx(vs[k-1 : k+2]...); min < v {
			mp[k-1+byte(n)] = min
		} else {
			mp[k] = v
		}
	}

	if nBottom < 0 {
		return mp
	}

	ks, vs = u8i64.Map2KVs(mp, nil, func(i, j int) bool { return i < j })
	mp = make(map[byte]int)
	for i := 0; i < nBottom && i < len(ks); i++ {
		mp[ks[i]] = vs[i]
	}
	return mp
}

// return positions of max-up-trend to max-down-trend
func Slope(data []float64, step int) (sp []int) {

	slope := []struct {
		ix int
		dy float64
	}{}

	for i := step - 1; i < len(data); i++ {
		a := data[i-(step-1)]
		b := data[i]
		slope = append(slope, struct {
			ix int
			dy float64
		}{
			ix: i - step/2,
			dy: b - a,
		})
	}

	sort.SliceStable(slope, func(i, j int) bool {
		return slope[i].dy > slope[j].dy
	})

	for _, s := range slope {
		sp = append(sp, s.ix)
	}
	return
}
