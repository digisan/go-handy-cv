package gocv

import (
	"fmt"
	"image"
	"image/color"
	"log"

	gio "github.com/digisan/gotk/io"
	cm "github.com/digisan/tiny-gocv/math"
)

const (
	ClrErr = 4
)

func Color2RGBA(clr color.Color) (r, g, b, a byte) {
	r32, g32, b32, a32 := clr.RGBA()
	return byte(r32 >> 8), byte(g32 >> 8), byte(b32 >> 8), byte(a32 >> 8)
}

func Color2Gray(clr color.Color) byte {
	gray32, _, _, _ := clr.RGBA()
	return byte(gray32 >> 8)
}

func ColorEqual(c1, c2 color.RGBA, eR, eG, eB byte) bool {
	// if cm.DisByte(c1.R, c2.R) < eR && cm.DisByte(c1.G, c2.G) < eG && cm.DisByte(c1.B, c2.B) < eB && cm.DisByte(c1.A, c2.A) < eA {
	// 	return true
	// }
	if cm.DisByte(c1.R, c2.R) < eR && cm.DisByte(c1.G, c2.G) < eG && cm.DisByte(c1.B, c2.B) < eB {
		return true
	}
	return false
}

func CompositeRGBA(r, g, b, a image.Image) *image.RGBA {

	rectR, rectG, rectB, rectA := r.Bounds(), g.Bounds(), b.Bounds(), a.Bounds()
	if rectR != rectG || rectG != rectB || rectB != rectA {
		log.Fatalln("r, g, b, a all must be same size")
		return nil
	}

	rgba := image.NewRGBA(rectR)
	bytes := rgba.Pix
	for i, p := range r.(*image.Gray).Pix {
		bytes[i*4] = p
	}
	for i, p := range g.(*image.Gray).Pix {
		bytes[i*4+1] = p
	}
	for i, p := range b.(*image.Gray).Pix {
		bytes[i*4+2] = p
	}
	for i, p := range a.(*image.Gray).Pix {
		bytes[i*4+3] = p
	}
	return rgba
}

func FindColorArea(img image.Image, clr color.RGBA) map[image.Point]struct{} {

	// out
	mPt := make(map[image.Point]struct{})

	r, g, b, _ := SplitRGBA(img)
	rect := img.Bounds()

	// order is same as below for range
	mChClr := map[int]byte{
		0: clr.R,
		1: clr.G,
		2: clr.B,
	}
	mPtN := make(map[image.Point]int)
	for i, ch := range []*image.Gray{r, g, b} {
		for y := 0; y < rect.Dy(); y++ {
			offset := y * ch.Stride
			line := ch.Pix[offset:]
			for x := 0; x < rect.Dx(); x++ {
				pxl := line[x]
				pt := image.Point{X: x, Y: y}
				if cm.DisByte(pxl, mChClr[i]) < ClrErr {
					mPtN[pt]++
				}
			}
		}
	}

	for pt, n := range mPtN {
		if n == 3 {
			mPt[pt] = struct{}{}
		}
	}

	return mPt
}

func SplitRGBA(img image.Image) (r, g, b, a *image.Gray) {

	rect := img.Bounds()

	left, top, right, bottom := rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y
	img = ROI4RGBA(img, left, top, right, bottom)

	var bytes []byte
	switch pImg := img.(type) {
	case *image.RGBA:
		bytes = pImg.Pix
	case *image.NRGBA:
		bytes = pImg.Pix
	// case *image.YCbCr: //	YCbCrSubsampleRatio444
	// 	bytes = pImg.Pix
	default:
		log.Fatalf("[%v] is not support", pImg)
	}

	r, g, b, a = image.NewGray(rect), image.NewGray(rect), image.NewGray(rect), image.NewGray(rect)
	for i, p := range bytes {
		switch i % 4 {
		case 0:
			r.Pix[i/4] = p
		case 1:
			g.Pix[i/4] = p
		case 2:
			b.Pix[i/4] = p
		case 3:
			a.Pix[i/4] = p
		}
	}

	// wg := &sync.WaitGroup{}
	// wg.Add(4)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[0:], r.Pix)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[1:], g.Pix)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[2:], b.Pix)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[3:], a.Pix)
	// wg.Wait()

	return
}

func FindPosByColor(img image.Image, c color.RGBA) (pos []image.Point) {
	rect := img.Bounds()
	rgba := ROI4RGBA(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
	for y := 0; y < rect.Dy(); y++ {
		start := y * rgba.Stride
		pln := rgba.Pix[start : start+rgba.Stride]
		for x := 0; x < rect.Dx(); x++ {
			p := pln[4*x:]
			cmp := color.RGBA{p[0], p[1], p[2], p[3]}
			if ColorEqual(c, cmp, 3, 3, 3) {
				pos = append(pos, image.Point{x, y})
			}
		}
	}
	return
}

func FindROIByColor(img image.Image, c color.RGBA, sRadius, iRadius float64, auditPath string) (mPtROI map[image.Point]*image.RGBA) {

	mPtRGBA := make(map[image.Point]*image.RGBA)
	for _, pos := range FindPosByColor(img, c) {
		roi := ROI4RGBAv2(img, pos.X, pos.Y, int(sRadius))
		mPtRGBA[pos] = roi
	}

	mPtROI = make(map[image.Point]*image.RGBA)
NEXT:
	for pt1, rgba := range mPtRGBA {
		for pt2 := range mPtROI {
			if cm.DisPt(pt1, pt2) < iRadius {
				continue NEXT
			}
		}
		mPtROI[pt1] = rgba
	}

	I := 0
	for pt, roi := range mPtROI {
		if len(auditPath) > 0 {
			gio.MustCreateDir(auditPath)
			SavePNG(roi, fmt.Sprintf("./%s/%00d-%d-%d.png", "./out/audit/", I, pt.X, pt.Y))
			I++
		}
	}

	return
}

// func FindROIrgbaByBlob(img image.Image,
// 	sRadius int,
// 	filterR func(x, y int, p byte) bool,
// 	filterG func(x, y int, p byte) bool,
// 	filterB func(x, y int, p byte) bool,
// 	disErr int,
// 	auditPath string) (mPtRGBA map[image.Point]*image.RGBA) {

// 	gotkio.MustCreateDir(auditPath)
// 	mPtRGBA = make(map[image.Point]*image.RGBA)

// 	// rect := img.Bounds()
// 	// rgba := ROIrgba(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)

// 	r, g, b, _ := SplitRGBA(img)
// 	blobPosGrp := gocv.DetectClrBlobPos(r.Rect.Dx(), r.Rect.Dy(), r.Stride,
// 		r.Pix, g.Pix, b.Pix,
// 		filterR, filterG, filterB, disErr)

// 	for i, bpos := range blobPosGrp {
// 		pos := image.Point{X: bpos.X, Y: bpos.Y}
// 		roi := ROIrgbaV2(img, pos.X, pos.Y, sRadius)
// 		mPtRGBA[pos] = roi
// 		if len(auditPath) > 0 {
// 			savePNG(roi, fmt.Sprintf("./%s/%00d-%d-%d.png", auditPath, i, pos.X, pos.Y))
// 		}
// 	}
// 	return
// }
