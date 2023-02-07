package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fajarardiyanto/steganography/vendor/bimg"
	"github.com/signintech/pdft"
)

const inputFile = "file/input.jpg"
const destTwo = "file/inject_two.jpg"

func main() {
	// ImgToPDF()
	removeMetadata()
}

func ImgToPDF() {
	var pt pdft.PDFt
	err := pt.Open(inputFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = pt.Save(destTwo)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func removeMetadata() {
	buffer, err := bimg.Read(inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	imgs := bimg.NewImage(buffer)
	size, _ := imgs.Size()
	imgt := imgs.Type()
	newImage, err := imgs.
		Process(bimg.Options{
			RemoveAllMetaData: true,
			StripMetadata:     true,
			Width:             size.Width,
			Height:            size.Height,
			Embed:             true,
			Quality:           50,
			Compression:       9,
			Interpolator:      bimg.Bilinear,
			Trim:              true,
			Type:              bimg.PNG,
			NoProfile:         true,
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	imgType := bimg.PNG
	switch imgt {
	case bimg.ImageTypeName(bimg.JPEG):
		imgType = bimg.JPEG
	case bimg.ImageTypeName(bimg.PDF):
		imgType = bimg.JPEG
	}

	dat, err := bimg.NewImage(newImage).Process(bimg.Options{
		Type:              imgType,
		RemoveAllMetaData: true,
		Trim:              true,
		StripMetadata:     true,
		NoProfile:         true,
		Embed:             true,
		Compression:       9, Quality: 50})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	if err = bimg.Write(destTwo, dat); err != nil {
		log.Printf("error write file image two: %v", err)
		return
	}
}
