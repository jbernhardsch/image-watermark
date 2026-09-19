package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	// CAM_CANON - for canon camera
	CAM_CANON = "canon"
	// CAM_FUJIFILM - for fujifilm camera
	CAM_FUJIFILM = "fujifilm"

	// IMAGE_NOEXIF - for all images with no exif
	IMAGE_NOEXIF = "no-exif"

	// IMAGE_LANDSCAPE - landscape image
	IMAGE_LANDSCAPE = "landscape"
	// IMAGE_PORTRAIT - portrait image
	IMAGE_PORTRAIT = "portrait"
	// IMAGE_PORTRAIT_FUJI - portrait image for fuji film
	IMAGE_PORTRAIT_FUJI = "portrait-fuji"
	// IMAGE_PORTRAIT_FUJI_SPC - portrait image for fuji film (special rotation)
	IMAGE_PORTRAIT_FUJI_SPC = "portrait-fuji-spc"

	// RESIZE_VALUE - default value for resizing
	RESIZE_VALUE = 2000
)

func main() {
	var dir = flag.String("dir", "", "type the directory of images")
	var waterMark = flag.String("watermark", "", "image for watermark")
	var size = flag.Int("size", RESIZE_VALUE, "resize target (long edge, in pixels)")
	flag.Parse()

	if *size > 0 {
		RESIZE_VALUE = *size
	}

	// reserve target directory name
	targetDir := fmt.Sprintf("%s/watermarked", *dir)

	// check the directory then create the watermarked directory
	if err := checkDirectory(*dir, targetDir); err != nil {
		log.Fatal(err)
	}

	if err := proceedDirectory(*dir, targetDir, *waterMark); err != nil {
		log.Fatal(err)
	}
}
