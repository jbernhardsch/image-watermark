package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// resizeImaging - function for resizing image
func resizeImaging(needRotate bool, orientation, fullName, dir, targetDir string) error {
	imgSrc := fmt.Sprintf("%s/%s", dir, fullName)
	input, err := imaging.Open(imgSrc)
	if err != nil {
		return fmt.Errorf("error open the image %s: %+v", fullName, err)
	}

	var (
		width, height int
	)
	// set image size based on orientation and camera
	if (orientation == IMAGE_LANDSCAPE && !needRotate) || ((orientation == IMAGE_PORTRAIT_FUJI || orientation == IMAGE_PORTRAIT_FUJI_SPC) && needRotate) || (orientation == IMAGE_LANDSCAPE && needRotate) {
		width = RESIZE_VALUE
		height = 0

	} else if orientation == IMAGE_PORTRAIT && needRotate {
		width = 0
		height = RESIZE_VALUE

	} else if orientation == IMAGE_NOEXIF {
		if reader, err := os.Open(filepath.Join(dir, fullName)); err == nil {
			defer reader.Close()
			im, _, _ := image.DecodeConfig(reader)

			if im.Width > im.Height {
				width = RESIZE_VALUE
				height = 0
			} else {
				width = 0
				height = RESIZE_VALUE
			}

		} else {
			fmt.Println("Impossible to open the file:", err)
		}

	} else {
		width = RESIZE_VALUE
		height = 0
	}

	input = imaging.Resize(input, width, height, imaging.Lanczos)

	// if need rotate image then it will rotate the image
	if needRotate {
		if orientation == IMAGE_PORTRAIT_FUJI_SPC {
			input = imaging.Rotate270(input)
		} else {
			input = imaging.Rotate90(input)
		}
	}

	// save to target directory
	err = imaging.Save(input, fmt.Sprintf("%s/%s", targetDir, fullName))
	if err != nil {
		return fmt.Errorf("failed to save the image %s: %+v", fullName, err)
	}

	return nil
}
