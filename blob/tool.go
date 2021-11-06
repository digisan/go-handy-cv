package blob

import (
	"fmt"
	"image"
	"math"
	"sort"
	"strings"
)

var (
	sTrimRight  = strings.TrimRight
	sSplit      = strings.Split
	sTrimSuffix = strings.TrimSuffix
	sIndex      = strings.Index
	sLastIndex  = strings.LastIndex
	sJoin       = strings.Join
	sContains   = strings.Contains
	sHasPrefix  = strings.HasPrefix
)

func PtDis(pt1, pt2 image.Point) int {
	dx := math.Abs(float64(pt1.X - pt2.X))
	dy := math.Abs(float64(pt1.Y - pt2.Y))
	return int(math.Sqrt(dx*dx + dy*dy))
}

func PtsAverage(pts ...image.Point) image.Point {
	sumX, sumY := 0, 0
	for _, pt := range pts {
		sumX += pt.X
		sumY += pt.Y
	}
	return image.Point{X: sumX / len(pts), Y: sumY / len(pts)}
}

func PtInRect(pt image.Point, rect image.Rectangle) bool {
	lefttop := rect.Min
	rightbottom := rect.Max
	return (pt.X > lefttop.X && pt.X < rightbottom.X) && (pt.Y > lefttop.Y && pt.Y < rightbottom.Y)
}

func RectCrossed(rect1, rect2 image.Rectangle) bool {
	lefttop1 := rect1.Min
	rightbottom1 := rect1.Max
	lefttop2 := rect2.Min
	rightbottom2 := rect2.Max
	if (lefttop1.X >= lefttop2.X && lefttop1.X <= rightbottom2.X) &&
		(rightbottom1.X >= lefttop2.X && rightbottom1.X <= rightbottom2.X) &&
		(lefttop1.Y <= lefttop2.Y && rightbottom1.Y >= rightbottom2.Y) {
		return true
	}
	return false
}

func RectOverlap(rect1, rect2 image.Rectangle) bool {
	for i, rect := range []image.Rectangle{rect1, rect2} {
		compare := rect2
		if i == 1 {
			compare = rect1
		}

		if RectCrossed(rect, compare) {
			return true
		}

		lefttop := rect.Min
		rightbottom := rect.Max
		leftbottom := image.Point{lefttop.X, rightbottom.Y}
		righttop := image.Point{rightbottom.X, lefttop.Y}
		for _, pt := range []image.Point{lefttop, rightbottom, leftbottom, righttop} {
			if PtInRect(pt, compare) {
				return true
			}
		}
	}
	return false
}

func Map2KVs4BL(m map[int][]*blobline, less4key func(i int, j int) bool) (keys []int, values [][]*blobline) {

	if m == nil {
		return nil, nil
	}
	if len(m) == 0 {
		return []int{}, [][]*blobline{}
	}

	type kv struct {
		key   int
		value []*blobline
	}

	kvSlc := []kv{}
	for k, v := range m {
		kvSlc = append(kvSlc, kv{key: k, value: v})
	}

	switch {
	case less4key != nil:
		sort.SliceStable(kvSlc, func(i, j int) bool { return less4key(kvSlc[i].key, kvSlc[j].key) })
	default:
		// do not sort
	}

	for _, kvEle := range kvSlc {
		keys = append(keys, kvEle.key)
		values = append(values, kvEle.value)
	}
	return
}

func ParseIntPair(pair string) (data [2]int) {
	fmt.Sscanf(pair, "[%d,%d]", &data[0], &data[1])
	return
}
