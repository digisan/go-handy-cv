package blob

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Blob struct {
	Y   int
	Idx int
	tag string
}

func (b Blob) Loc() [2]Point {
	top, bottom := -1, -1
	start, end := 1000000, -1000000

	taglns := sSplit(b.tag, "\n")
	for i, ltag := range taglns {
		yse := sSplit(ltag, ":")
		//
		yStr := yse[0]
		switch i {
		case 0:
			top, _ = strconv.Atoi(yStr)
			top--
		case len(taglns) - 1:
			bottom, _ = strconv.Atoi(yStr)
			bottom++
		}
		//
		seStr := sSplit(sTrim(sTrim(yse[1], " "), "[]"), ",")
		sStr, eStr := seStr[0], seStr[1]
		if s, _ := strconv.Atoi(sStr); s < start {
			start = s - 1
		}
		if e, _ := strconv.Atoi(eStr); e > end {
			end = e
		}
	}
	return [2]Point{{start, top}, {end, bottom}}
}

func (b Blob) Area() (area int) {
	taglns := sSplit(b.tag, "\n")
	for _, ltag := range taglns {
		yse := sSplit(ltag, ":")
		seStr := sSplit(sTrim(sTrim(yse[1], " "), "[]"), ",")
		sStr, eStr := seStr[0], seStr[1]
		s, _ := strconv.Atoi(sStr)
		e, _ := strconv.Atoi(eStr)
		area += (e - s)
	}
	return area
}

func MkSet(blobs ...Blob) (rtBlobs []Blob) {
	m := make(map[string]Blob)
	for _, blob := range blobs {
		m[blob.tag] = blob
	}
	for _, blob := range m {
		rtBlobs = append(rtBlobs, blob)
	}
	return
}

///

type blobline struct {
	y          int
	start, end int
	prev, next []*blobline
}

func (bl blobline) String() string {
	return fmt.Sprintf("%d: [%d,%d]", bl.y, bl.start, bl.end)
}

func traverse(bl *blobline, sb *strings.Builder, f func(bl *blobline, sb *strings.Builder)) {

	f(bl, sb)

	if len(bl.next) == 0 {
		return
	}
	for _, next := range bl.next {
		traverse(next, sb, f)
	}
}

func adjacent(a1, a2 [2]int) bool {
	start1, end1 := a1[0], a1[1]
	start2, end2 := a2[0], a2[1]
	switch {
	// ***e
	//    ***e
	case end1 <= start2:
		return false
	//    ***e
	// ***e
	case end2 <= start1:
		return false
	default:
		return true
	}
}

// return the first above blobline
func blAdjacent(bl1, bl2 *blobline) (*blobline, bool) {
	if bl2.y-bl1.y == 1 {
		if adjacent([2]int{bl1.start, bl1.end}, [2]int{bl2.start, bl2.end}) {
			bl1.next = append(bl1.next, bl2)
			bl2.prev = append(bl2.prev, bl1)
			return bl1, true
		}
	}
	return bl1, false
}

// blob line scan
func scan(y int, line []byte) (lob []*blobline) {
	inBlob := false
	bl := &blobline{y: y}
	for i, p := range line {
		if !inBlob && p == 0 {
			inBlob = true
			bl.start = i
			bl.end = len(line)
		}
		if inBlob && p > 0 {
			inBlob = false
			bl.end = i
			lob = append(lob, bl)
			bl = &blobline{y: y} // create a new one after appending
		}
	}
	if inBlob {
		lob = append(lob, bl) // blob to edge
	}
	return
}

func DetectBlob(width, height, step int, data []byte) []Blob {

	mBlobsLine := make(map[int][]*blobline)

	for y := 0; y < height; y++ {
		yIdx := y * step
		line := data[yIdx : yIdx+width]
		mBlobsLine[y] = scan(y, line)

		// for i := 0; i < width; i++ {
		// 	idx := yIdx + i
		// 	p := &data[idx]
		// 	*p *= 2
		// }
	}

	keys, values := Map2KVs4BL(mBlobsLine, func(i, j int) bool { return i < j })

	// judge connection by each adjacent blob line pairs

	mYBlobs := make(map[int][]*blobline)

	for i := 0; i < len(keys)-1; i++ {

		yThis, yBelow := keys[i], keys[i+1]
		blsThis, blsBelow := values[yThis], values[yBelow]

		for _, blThis := range blsThis {
			linked := false
			for _, blBelow := range blsBelow {
				if bl, ok := blAdjacent(blThis, blBelow); ok {
					if len(bl.prev) == 0 {
						mYBlobs[yThis] = append(mYBlobs[yThis], bl)
						linked = true
					}
				}
			}
			if !linked {
				if len(blThis.prev) == 0 && len(blThis.next) == 0 {
					mYBlobs[yThis] = append(mYBlobs[yThis], blThis)
				}
			}
		}
	}
	{
		yBottom := keys[len(keys)-1]
		blsBottom := values[yBottom]
		for _, blBottom := range blsBottom {
			if len(blBottom.prev) == 0 {
				mYBlobs[yBottom] = append(mYBlobs[yBottom], blBottom)
			}
		}
	}

	// -----------------------------------------------

	mYWriters := make(map[int][]*strings.Builder)

	f := func(blob *blobline, sb *strings.Builder) {
		sb.WriteString(fmt.Sprintln(blob))
	}

	keys, values = Map2KVs4Blob(mYBlobs, func(i, j int) bool { return i < j })
	for i, y := range keys {
		blobs := values[i]
		for _, blob := range blobs {
			sb := &strings.Builder{}
			mYWriters[y] = append(mYWriters[y], sb)
			traverse(blob, sb, f)
		}
	}

	blobs := []Blob{}
	for y, writers := range mYWriters {
		for idx, w := range writers {
			blob := Blob{Y: y, Idx: idx, tag: sTrimRight(w.String(), "\n")}
			blobs = append(blobs, blob)
		}
	}
	return MkSet(blobs...) // remove duplicated blob
}
