package main

import (
	"fmt"
	"os"
	"strings"
)

// checkDirectory - function to check whether directory exists or not
func checkDirectory(dir, targetDir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("the main folder %s doesn't exist", dir)
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		// if exists then create watermarked directory
		if err := os.Mkdir(targetDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create watermarked folder")
		}
	}

	return nil
}

// proceedDirectory - function for proceeding the directory
func proceedDirectory(dir, targetDir, wmImage string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %+v", dir, err)
	}

	for _, file := range files {
		// lower string file name
		src := strings.ToLower(file.Name())

		fName := strings.Split(src, ".")
		if fName[len(fName)-1] == "jpg" {
			fmt.Println("processing image", file.Name()) // TO DO: remove it

			if err := proceedImaging(file.Name(), wmImage, dir, targetDir); err != nil {
				//if err := proceedImage(file.Name(), fName[0], dir, targetDir); err != nil {
				return fmt.Errorf("failed to proceed %s: %+v", file.Name(), err)
			}
		}
	}

	return nil
}

// proceedImaging - function for imaging process
func proceedImaging(fullName, wmImage, dir, targetDir string) error {
	// check orientation from its exif
	ori, needRotate, err := getOrientation(fullName, dir)
	if err != nil {
		return fmt.Errorf("cannot get the orientation %s:%+v", fullName, err)
	}

	if err := resizeImaging(needRotate, ori, fullName, dir, targetDir); err != nil {
		return fmt.Errorf("cannot resize image %s:%+v", fullName, err)
	}

	if err := waterMarking(fullName, wmImage, targetDir); err != nil {
		return fmt.Errorf("cannot add watermark to %s:%+v", fullName, err)
	}

	return nil
}
