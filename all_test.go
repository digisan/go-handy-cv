package gocv

import (
	"fmt"
	"image"
	"testing"
)

func TestRangePixelMap(t *testing.T) {

	img := LoadImage("./data/sample.jpg")
	// img := loadImage("./cfg/roi.jpg")

	gray := ToGray(img)
	fmt.Println("Cvt2Gray Done")

	red, green, blue, alpha := SplitRGBA(img)
	fmt.Println("Cvt2RGBA Done")

	edgeCentres := []image.Point{}

	for x := 1800; x < 2800; x += 50 {

		for y := 1100; y > 200; y -= 20 {

			centre := image.Point{X: x, Y: y}
			pts := RoundPts(centre, 100, 5)

			// for _, pt := range pts {
			// 	gio.MustAppendFile("pts.txt", []byte(fmt.Sprint(pt.X, pt.Y)), true)
			// }

			ppts := Pts2PPts(pts, gray, red, green, blue, alpha, centre)

			N := len(ppts)

			I := 0
			// plptsOut := []PolarPoint{}
			for _, pt := range ppts {
				// if pt.Phi <= 0 {
				if pt.Red > 100 && pt.Green > 100 {
					// plptsOut = append(plptsOut, pt)
					I++
				}
				// }
			}

			// pts = PlPtsXY(plptsOut)

			ratio := float64(I) / float64(N)
			if ratio > 0.25 && ratio < 0.35 {
				// DrawCircles(img, pts, 1, "G", "./out/output.jpg")
				// DrawCircle(img, centre, 10, "G", fmt.Sprintf("./out/%d-%d.jpg", x, y))
				edgeCentres = append(edgeCentres, centre)
				break
			}

			// ptsval := PtsVal(gray, pts)
			// fmt.Println("PtsVal Done", len(ptsval))

			// m, _, _ := histogram(ptsval)
			// imgHis := DrawHisto(m, Peaks(m, 3, 1, 2), nil)
			// saveJPG(imgHis, "./his1.jpg")

		}
	}

	DrawCircles(img, edgeCentres, 5, "G", "./out/edge.jpg")

	DrawSpline(img, edgeCentres, 5, "G", "./out/edge1.jpg")
}
