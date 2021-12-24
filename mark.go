package gocv

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
)

func ColorAreaToJSON(markedImg, outJSON, clr string) {

	var r, g, b byte
	switch clr {
	case "W", "w", "White", "white", "WHITE":
		r, g, b = 255, 255, 255
	case "K", "k", "Black", "black", "BLACK":
		r, g, b = 0, 0, 0
	case "R", "r", "Red", "red", "RED":
		r, g, b = 255, 0, 0
	case "G", "g", "Green", "green", "GREEN":
		r, g, b = 0, 255, 0
	case "B", "b", "Blue", "blue", "BLUE":
		r, g, b = 0, 0, 255
	case "C", "c", "Cyan", "cyan", "CYAN":
		r, g, b = 0, 255, 255
	case "Y", "y", "Yellow", "yellow", "YELLOW":
		r, g, b = 255, 255, 0
	case "M", "m", "Magenta", "magenta", "MAGENTA":
		r, g, b = 255, 0, 255
	default:
		r, g, b = 0, 0, 0
	}

	// search
	img := LoadImage(markedImg)
	mPt := FindColorArea(img, color.RGBA{r, g, b, 0})
	fmt.Println(len(mPt))

	// store
	mPtStore := map[string]interface{}{}
	for pt := range mPt {
		mPtStore[pt.String()] = ""
	}
	data, err := json.Marshal(mPtStore)
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outJSON, data, os.ModePerm)
}

func loadAreaFromJSON(areaJSON string) map[image.Point]struct{} {

	// load
	data, err := os.ReadFile(areaJSON)
	if err != nil {
		log.Fatalln(err)
	}
	mPtLoad := map[string]interface{}{}
	err = json.Unmarshal(data, &mPtLoad)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(len(mPtLoad))

	//////////////////////////////////

	mPt := make(map[image.Point]struct{})
	for strPt := range mPtLoad {
		pt := image.Point{}
		fmt.Sscanf(strPt, "(%d,%d)", &pt.X, &pt.Y)
		mPt[pt] = struct{}{}
	}

	return mPt
}

func LoadAreaFromJSON(areaJSON string, offsetX, offsetY int) (pts []image.Point) {
	for pt := range loadAreaFromJSON(areaJSON) {
		pts = append(pts, image.Pt(pt.X+offsetX, pt.Y+offsetY))
	}
	return
}

func PaintArea(imgFile, outImgFile, color string, pts []image.Point) {
	DrawCircles(LoadImage(imgFile), pts, 1, color, outImgFile)
}

func PaintAreaFromJSON(imgFile, areaJSON string, offsetX, offsetY int, outImgFile, color string) {
	// load & re-mark
	pts := LoadAreaFromJSON(areaJSON, offsetX, offsetY)
	PaintArea(imgFile, outImgFile, color, pts)
}
