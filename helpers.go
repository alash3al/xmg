package main

import (
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"

	"github.com/disintegration/imaging"
)

// processSingleFileUpload ...
func processSingleFileUpload(file *multipart.FileHeader) (image.Image, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dstfilename := filepath.Join(os.TempDir(), file.Filename)
	dst, err := os.Create(dstfilename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	return imaging.Open(dstfilename)
}

// processImageHash ...
func processImageHash(img image.Image) (*goimagehash.ImageHash, error) {
	img = imaging.Grayscale(img)
	return goimagehash.PerceptionHash(img)
}

// processImagesHashes ...
func processImagesHashes(imgs []image.Image) []*goimagehash.ImageHash {
	hashes := []*goimagehash.ImageHash{}
	for _, img := range imgs {
		hash, err := processImageHash(img)
		if err != nil {
			continue
		}

		hashes = append(hashes, hash)
	}

	return hashes
}

func getAllImageOrientations(img image.Image) []image.Image {
	ret := []image.Image{}

	for i := 0; i < 3; i++ {
		if i > 0 {
			img = ret[i-1]
		}
		ret = append(ret, imaging.Rotate90(img))
	}

	for _, img := range ret {
		ret = append(ret, imaging.FlipH(img), imaging.FlipV(img))
	}

	return ret
}
