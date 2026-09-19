package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

// getOrientation - function for checking orientation
func getOrientation(fullName, dir string) (string, bool, error) {
	input, err := os.Open(fmt.Sprintf("%s/%s", dir, fullName))
	if err != nil {
		return "", false, fmt.Errorf("failed to open file %s: %+v", fullName, err)
	}
	defer input.Close()

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	ex, err := exif.Decode(input)
	if err != nil {
		return "", false, fmt.Errorf("failed to decode image: %+v", err)
	}

	maker, err := ex.Get(exif.Make)
	if err != nil {
		// skip the process to directly check the size and resize
		return IMAGE_NOEXIF, false, nil

		//return "", false, fmt.Errorf("failed to get maker: %+v", err)
	}

	cam := strings.ReplaceAll(strings.ToLower(maker.String()), `"`, ``)

	// get orientation for images having exif
	ori, _ := ex.Get(exif.Orientation)

	// differentiate process based on camera brand
	if cam == CAM_CANON {
		// get length of x dimension
		xLength, err := ex.Get(exif.PixelXDimension)
		if err != nil {
			return "", false, fmt.Errorf("failed to get x dimension: %+v", err)
		}
		xStr := xLength.String()
		xSize, err := strconv.Atoi(xStr)
		if err != nil {
			return "", false, fmt.Errorf("failed to get x dimension: %+v", err)
		}

		// get length of y dimension
		yLength, err := ex.Get(exif.PixelYDimension)
		if err != nil {
			return "", false, fmt.Errorf("failed to get y dimension: %+v", err)
		}
		yStr := yLength.String()
		ySize, err := strconv.Atoi(yStr)
		if err != nil {
			return "", false, fmt.Errorf("failed to get y dimension: %+v", err)
		}

		// special action for portrait image
		if ori.String() == "8" && xSize > ySize {
			// rotate image first
			return IMAGE_LANDSCAPE, true, nil
		} else if ori.String() == "8" && xSize < ySize {
			// rotate image first
			return IMAGE_PORTRAIT, true, nil
		}

		// compare x and y dimension, x bigger means landscape
		if xSize > ySize {
			return IMAGE_LANDSCAPE, false, nil
		} else {
			return IMAGE_PORTRAIT, true, nil
		}

	} else if cam == CAM_FUJIFILM {
		if ori.String() == "1" {
			return IMAGE_LANDSCAPE, false, nil

		} else if ori.String() == "8" {
			// rotate image first
			return IMAGE_PORTRAIT_FUJI, true, nil

		} else {
			return IMAGE_PORTRAIT_FUJI_SPC, true, nil

		}
	}

	return "", false, errors.New("cannot proceed image orientation")
}
