package gocv

import (
	"image"
	"log"
	"os"
)

func LoadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(fmtName)

	return img
}
