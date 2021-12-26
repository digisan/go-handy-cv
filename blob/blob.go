package blob

import (
	"crypto/md5"
	"fmt"
	"image"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/digisan/go-generics/i64"
	"github.com/digisan/go-generics/str"
)

type Blob struct {
	y   int
	idx string
	tag string
	loc image.Rectangle
}

func tagsort(tag string) string {

	tln := sSplit(tag, "\n")
	str.Filter(&tln, func(i int, e string) bool { return sContains(e, ":") })
	if len(tln) > 1 {
		sort.Slice(tln, func(i, j int) bool {
			iln, jln := tln[i], tln[j]
			pi, pj := sIndex(iln, ":"), sIndex(jln, ":")
			ni, _ := strconv.Atoi(iln[:pi])
			nj, _ := strconv.Atoi(jln[:pj])
			if ni != nj {
				return ni < nj
			} else {
				pis, pjs := sIndex(iln, "["), sIndex(jln, "[")
				pie, pje := sIndex(iln, ","), sIndex(jln, ",")
				ni, _ := strconv.Atoi(iln[pis+1 : pie])
				nj, _ := strconv.Atoi(jln[pjs+1 : pje])
				return ni < nj
			}
		})
	}

	// combine
	s := 0
AGAIN:
	for i := s; i < len(tln)-1; i++ {
		p := sIndex(tln[i], ":")
		pfx := tln[i][:p+1]
		if strings.HasPrefix(tln[i+1], pfx) {
			tln[i] += tln[i+1][p+1:]
			tln = append(tln[:i+1], tln[i+2:]...)
			s = i
			goto AGAIN
		}
	}

	for i := 0; i < len(tln); i++ {
		p := sIndex(tln[i], ":") + 2
		pfx := tln[i][:p]
		pairs := sSplit(tln[i][p:], " ")
		pairs = str.MkSet(pairs...)
		if len(pairs) > 1 {
			sort.Slice(pairs, func(i, j int) bool {
				pi, pj := pairs[i], pairs[j]
				pis, pjs := sIndex(pi, "["), sIndex(pj, "[")
				pie, pje := sIndex(pi, ","), sIndex(pj, ",")
				ni, _ := strconv.Atoi(pi[pis+1 : pie])
				nj, _ := strconv.Atoi(pj[pjs+1 : pje])
				return ni < nj
			})
		}
		tln[i] = pfx + sJoin(pairs, " ")
	}

	return sJoin(tln, "\n")
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

func (b *Blob) Loc() image.Rectangle {

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

	b.loc = image.Rectangle{
		Min: image.Point{left, top},
		Max: image.Point{right, bottom},
	}
	return b.loc
}

func (b *Blob) Center() image.Point {
	if b.loc == image.Rect(0, 0, 0, 0) {
		b.Loc()
	}
	return image.Point{(b.loc.Max.X + b.loc.Min.X) / 2, (b.loc.Max.Y + b.loc.Min.Y) / 2}
}

func (b *Blob) Area() (area int) {
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
func scan(y int, line []byte, filter func(x, y int, p byte) bool) (lob []*blobline) {
	inBlob := false
	bl := &blobline{y: y}
	for i, p := range line {
		if !inBlob && filter(i, y, p) {
			inBlob = true
			bl.start = i
			bl.end = len(line)
		}
		if inBlob && !filter(i, y, p) {
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

func detectParts(width, height, step int, data []byte, filter func(x, y int, p byte) bool) []Blob {

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
			blob := Blob{y: y, idx: fmt.Sprint(idx), tag: tagsort(w.String())}
			blobs = append(blobs, blob)
		}
	}
	return mkSet(blobs...) // remove duplicated blob
}

func getY(tag string) int {
	if p := sIndex(tag, ":"); p >= 0 {
		ns := tag[:p]
		n, err := strconv.Atoi(ns)
		if err != nil {
			log.Fatalln(err)
		}
		return n
	}
	return -1
}

func getPairsByY(tag string, y int) string {
	for _, ln := range sSplit(tag, "\n") {
		pfx := fmt.Sprintf("%d: ", y)
		if sHasPrefix(ln, pfx) {
			return ln[len(pfx):]
		}
	}
	return ""
}

func pairlink(pair1, pair2 string) bool {
	if pair1 == "" || pair2 == "" {
		return false
	}
	pa1 := sSplit(pair1, " ")
	pa2 := sSplit(pair2, " ")
	for _, p1 := range pa1 {
		for _, p2 := range pa2 {
			if p1 == p2 {
				return true
			}
		}
	}
	return false
}

func pairmerge(pair1, pair2 string) string {
	if pair1 == pair2 {
		return pair1
	}
	pair12 := pair1 + " " + pair2
	pair12 = sTrimRight(pair12, " ")
	pairs := sSplit(pair12, " ")
	pairs = str.MkSet(pairs...)
	if len(pairs) >= 2 && pairs[0] != "" && pairs[1] != "" {
		sort.Slice(pairs, func(i, j int) bool {
			pis, pjs := sIndex(pairs[i], "["), sIndex(pairs[j], "[")
			pie, pje := sIndex(pairs[i], ","), sIndex(pairs[j], ",")
			ni, _ := strconv.Atoi(pairs[i][pis+1 : pie])
			nj, _ := strconv.Atoi(pairs[j][pjs+1 : pje])
			return ni < nj
		})
		return sJoin(pairs, " ")
	}
	if pairs[0] != "" {
		return pairs[0]
	}
	return pairs[1]
}

func getMaxMinY(tag string) (Ymin, Ymax int) {
	Ymin, Ymax = -1, -1

	if p := sIndex(tag, ":"); p >= 0 {
		ns := tag[:p]
		n, err := strconv.Atoi(ns)
		if err != nil {
			log.Fatalln(err)
		}
		Ymin = n
	}

	if p := sLastIndex(tag, ":"); p >= 0 {
		temp := tag[:p]
		if p1 := sLastIndex(temp, "\n"); p >= 0 {
			ns := temp[p1+1:]
			n, err := strconv.Atoi(ns)
			if err != nil {
				log.Fatalln(err)
			}
			Ymax = n
		}
	}

	switch {
	case Ymin >= 0 && Ymax >= 0:
		if Ymin > Ymax {
			Ymin, Ymax = Ymax, Ymin
		}
	case Ymin >= 0 && Ymax == -1:
		Ymax = Ymin
	case Ymax >= 0 && Ymin == -1:
		Ymin = Ymax
	}

	return
}

func merge2Blob(be1, be2 Blob) (merged Blob, shared bool) {

	minY1, maxY1 := getMaxMinY(be1.tag)
	minY2, maxY2 := getMaxMinY(be2.tag)

	pairsarr1, pairsarr2 := []string{}, []string{}
	ys := []int{}

	minY, maxY := i64.Min(minY1, minY2), i64.Max(maxY1, maxY2)
	for y := minY; y <= maxY; y++ {
		pairs1 := getPairsByY(be1.tag, y)
		pairsarr1 = append(pairsarr1, pairs1)
		pairs2 := getPairsByY(be2.tag, y)
		pairsarr2 = append(pairsarr2, pairs2)
		ys = append(ys, y)
		if !shared {
			shared = pairlink(pairs1, pairs2)
		}
	}

	mergedTag := []string{}

	if shared {
		for i := 0; i < len(ys); i++ {
			p1, p2 := pairsarr1[i], pairsarr2[i]
			pair := pairmerge(p1, p2)
			mergedTag = append(mergedTag, fmt.Sprintf("%d: %s", ys[i], pair))
		}
	}

	tag := sJoin(mergedTag, "\n")
	return Blob{y: getY(tag), tag: tagsort(tag), idx: "merged"}, shared
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
		if str.In(b.ID(), ids...) {
			ret = append(ret, b)
		}
	}
	return
}

func combine(blobs ...Blob) (ret []Blob) {

AGAIN:

	ret = []Blob{}
	merged := make(map[string]struct{})
	m := make(map[string][]string)
	linked := false

	for _, be1 := range blobs {
		linked = false
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

			if RectOverlap(be1.Loc(), be2.Loc()) {
				if _, ok := merge2Blob(be1, be2); ok {
					m[id1] = append(m[id1], id2)
					merged[id1] = struct{}{}
					merged[id2] = struct{}{}
					linked = true
				}
			}
		}
		if !linked {
			ret = append(ret, be1)
		}
	}

	for hID, tIDs := range m {
		toMerge := blobEGrp(blobs, append([]string{hID}, tIDs...)...)
		ret = append(ret, mergeToOneBlob(toMerge...))
	}

	if linked {
		blobs = ret
		goto AGAIN
	}

	return
}

func DetectBlob(width, height, step int, data []byte, filter func(x, y int, p byte) bool) []Blob {
	blobs := detectParts(width, height, step, data, filter)
	sort.Slice(blobs, func(i, j int) bool {
		return blobs[i].y < blobs[j].y
	})
	return combine(blobs...)
}

func DetectClrBlobPos(
	width, height, step int,
	dataR, dataG, dataB []byte,
	filterR, filterG, filterB func(x, y int, p byte) bool,
	disErr int) (pos []image.Point) {

	blobsR := DetectBlob(width, height, step, dataR, filterR)
	blobsG := DetectBlob(width, height, step, dataG, filterG)
	blobsB := DetectBlob(width, height, step, dataB, filterB)

	for _, bR := range blobsR {
		cR := bR.Center()
		for _, bG := range blobsG {
			cG := bG.Center()
			for _, bB := range blobsB {
				cB := bB.Center()
				if PtDis(cR, cG) <= disErr && PtDis(cR, cB) <= disErr && PtDis(cG, cB) <= disErr {
					pos = append(pos, PtsAverage(cR, cG, cB))
				}
			}
		}
	}
	return
}
