package gocv

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func SaveJPG(img image.Image, path string) image.Image {
	out, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100
	if err := jpeg.Encode(out, img, &opts); err != nil {
		log.Println("Error on Saving as JPG")
	}
	return img
}

func SavePNG(img image.Image, path string) image.Image {
	out, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		log.Println("Error on Saving as PNG")
	}
	return img
}
