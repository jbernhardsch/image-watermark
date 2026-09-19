package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
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

	// Read the true raster dimensions straight from the pixel data instead
	// of trusting the EXIF PixelXDimension/PixelYDimension tags. Editors
	// such as Photoshop bake the rotation into the pixels themselves and
	// reset (or drop) the Orientation tag, but they don't reliably keep the
	// EXIF dimension tags in sync with the edited pixel data - decoding the
	// file directly is the only source that can't go stale.
	width, height, err := decodedSize(fullName, dir)
	if err != nil {
		return "", false, fmt.Errorf("failed to read image dimensions: %+v", err)
	}

	// differentiate process based on camera brand
	if cam == CAM_CANON {
		// special action for portrait image
		if ori.String() == "8" && width > height {
			// rotate image first
			return IMAGE_LANDSCAPE, true, nil
		} else if ori.String() == "8" && width < height {
			// rotate image first
			return IMAGE_PORTRAIT, true, nil
		}

		// compare x and y dimension, x bigger means landscape
		if width > height {
			return IMAGE_LANDSCAPE, false, nil
		}

		// The raster is already portrait shaped even though the
		// orientation tag doesn't say "8" (e.g. it's "1", or missing).
		// Straight off a Canon body this never happens - the sensor
		// always records landscape-shaped pixels and relies on the "8"
		// tag above to signal rotation. But an editor like Photoshop
		// bakes the rotation into the pixels and resets the tag, so by
		// the time we get here the image is already correctly oriented
		// and must not be rotated again.
		return IMAGE_PORTRAIT, false, nil

	} else if cam == CAM_FUJIFILM {
		if width < height {
			// Same reasoning as the Canon branch above: pixels are
			// already portrait shaped (most likely an editor already
			// baked in the rotation), so there's nothing left to rotate.
			return IMAGE_PORTRAIT_FUJI, false, nil
		}

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

// decodedSize - reads the actual pixel dimensions straight from the image
// data, independent of any (possibly stale or edited) EXIF metadata
func decodedSize(fullName, dir string) (width, height int, err error) {
	f, err := os.Open(filepath.Join(dir, fullName))
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}
