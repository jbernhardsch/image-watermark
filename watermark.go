package main

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

// waterMarking - function for adding watermark
func waterMarking(fullName, wmImage, targetDir string) error {
	imgSrc := fmt.Sprintf("%s/%s", targetDir, fullName)
	mainImage, err := imaging.Open(imgSrc)
	if err != nil {
		return fmt.Errorf("error open the image %s: %+v", fullName, err)
	}

	markImage, err := imaging.Open(wmImage)
	if err != nil {
		return fmt.Errorf("error open the watermark image %s: %+v", fullName, err)
	}

	mainBound := mainImage.Bounds()
	markBound := markImage.Bounds()

	offset := image.Pt(
		(mainBound.Size().X/2)-(markBound.Size().X/2),
		(mainBound.Size().Y - (markBound.Size().Y + 80)))

	waterMarked := imaging.Overlay(mainImage, markImage, offset, 0.5)
	// save to target directory
	err = imaging.Save(waterMarked, fmt.Sprintf("%s/%s", targetDir, fullName))
	if err != nil {
		return fmt.Errorf("failed to save the image %s: %+v", fullName, err)
	}
	return nil
}
