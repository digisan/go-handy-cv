package blob

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/digisan/gotk/slice/ts"
)

type Point struct {
	X int
	Y int
}

type Blob struct {
	y   int
	idx string
	tag string
}

func (b Blob) String() string {
	sb := strings.Builder{}
	sb.WriteString("\n------------------------------------\n")
	sb.WriteString(fmt.Sprintf("ID: %s\n", b.ID()))
	sb.WriteString(fmt.Sprintf("Y: %d Index: %s\n", b.y, b.idx))
	sb.WriteString(b.Tag() + "\n")
	sb.WriteString(fmt.Sprintf("location: %v\n", b.Loc()))
	sb.WriteString(fmt.Sprintf("area: %d\n", b.Area()))
	return sb.String()
}

func (b Blob) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(b.tag)))
}

func (b Blob) Tag() string {
	return b.tag
}

func (b Blob) Loc() [2]Point {

	top, bottom := 8192, 0
	left, right := 8192, 0

	r := regexp.MustCompile(`\d+:`)
	ns := r.FindAllString(b.tag, -1)
	for _, n := range ns {
		n = sTrimSuffix(n, ":")
		num, _ := strconv.Atoi(n)
		if num < top {
			top = num
		}
		if num > bottom {
			bottom = num
		}
	}

	r = regexp.MustCompile(`\[\d+,\d+\]`)
	ns = r.FindAllString(b.tag, -1)
	for _, n := range ns {
		pair := ParseIntPair(n)
		if pair[0] < left {
			left = pair[0]
		}
		if pair[1] > right {
			right = pair[1] - 1
		}
	}

	return [2]Point{{left, top}, {right, bottom}}
}

func (b Blob) Area() (area int) {
	r := regexp.MustCompile(`\[\d+,\d+\]`)
	ns := r.FindAllString(b.tag, -1)
	for _, n := range ns {
		pair := ParseIntPair(n)
		area += (pair[1] - pair[0])
	}
	return area
}

func mkSet(blobs ...Blob) (rt []Blob) {
	m := make(map[string]Blob)
	for _, blob := range blobs {
		m[blob.tag] = blob
	}
	for _, blob := range m {
		rt = append(rt, blob)
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
func scan(y int, line []byte, filter func(p byte) bool) (lob []*blobline) {
	inBlob := false
	bl := &blobline{y: y}
	for i, p := range line {
		if !inBlob && filter(p) {
			inBlob = true
			bl.start = i
			bl.end = len(line)
		}
		if inBlob && !filter(p) {
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

func detectParts(width, height, step int, data []byte, filter func(p byte) bool) []Blob {

	mBlobsLine := make(map[int][]*blobline)

	for y := 0; y < height; y++ {
		yIdx := y * step
		line := data[yIdx : yIdx+width]
		mBlobsLine[y] = scan(y, line, filter)

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

	keys, values = Map2KVs4BL(mYBlobs, func(i, j int) bool { return i < j })
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
			blob := Blob{y: y, idx: fmt.Sprint(idx), tag: sTrimRight(w.String(), "\n")}
			blobs = append(blobs, blob)
		}
	}
	return mkSet(blobs...) // remove duplicated blob
}

func merge2Blob(be1, be2 Blob) (merged Blob, shared bool) {
	tagrln1 := ts.Reverse(sSplit(be1.tag, "\n"))
	tagrln2 := ts.Reverse(sSplit(be2.tag, "\n"))
	mergedTag := []string{}
	i := -1
	for {
		i++
		if i < len(tagrln1) && i < len(tagrln2) {
			if tagrln1[i] == tagrln2[i] {
				mergedTag = append(mergedTag, tagrln1[i])
				shared = true
				continue
			}
		}
		if !shared {
			return Blob{}, false
		}
		if shared && i < len(tagrln1) && i < len(tagrln2) {
			p := sIndex(tagrln2[i], ":") + 1
			mergedTag = append(mergedTag, tagrln1[i]+tagrln2[i][p:])
			continue
		}
		if i < len(tagrln1) {
			mergedTag = append(mergedTag, tagrln1[i])
			continue
		}
		if i < len(tagrln2) {
			mergedTag = append(mergedTag, tagrln2[i])
			continue
		}
		break
	}
	mergedTag = ts.Reverse(mergedTag)
	tag := sJoin(mergedTag, "\n")
	y, err := strconv.Atoi(tag[:sIndex(tag, ":")])
	if err != nil {
		panic(err)
	}
	return Blob{y: y, tag: tag, idx: "merged"}, shared
}

func mergeToOneBlob(blobs ...Blob) Blob {
	b := blobs[0]
	for i := 1; i < len(blobs); i++ {
		b, _ = merge2Blob(b, blobs[i])
	}
	return b
}

func blobEGrp(blobs []Blob, ids ...string) (ret []Blob) {
	for _, b := range blobs {
		if ts.In(b.ID(), ids...) {
			ret = append(ret, b)
		}
	}
	return
}

func combine(blobs ...Blob) (ret []Blob) {

	merged := make(map[string]struct{})
	m := make(map[string][]string)
	for _, be1 := range blobs {
		linked := false
		id1 := be1.ID()
		if _, ok := merged[id1]; ok {
			continue
		}
		for _, be2 := range blobs {
			id2 := be2.ID()
			if _, ok := merged[id2]; ok {
				continue
			}
			if id1 == id2 {
				continue
			}
			if _, ok := merge2Blob(be1, be2); ok {
				m[id1] = append(m[id1], id2)
				merged[id1] = struct{}{}
				merged[id2] = struct{}{}
				linked = true
			}
		}
		if !linked {
			ret = append(ret, be1)
		}
	}

	for hID, tIDs := range m {
		ids := append([]string{hID}, tIDs...)
		toMerge := blobEGrp(blobs, ids...)
		ret = append(ret, mergeToOneBlob(toMerge...))
	}

	return
}

func DetectBlob(width, height, step int, data []byte, filter func(p byte) bool) []Blob {
	return combine(detectParts(width, height, step, data, filter)...)
}
