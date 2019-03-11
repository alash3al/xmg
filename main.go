package main

import (
	"flag"
	"image"
	"log"
	"strings"

	"github.com/Kagami/go-face"

	"github.com/labstack/echo"
)

var (
	flagHTTPAddr      = flag.String("listen", ":4068", "the http listen address")
	flagDLIBModelsDir = flag.String("dlib-models", "./", "full path to the dlib models directory")
)

var (
	recognizer *face.Recognizer
)

func init() {
	flag.Parse()

	rec, err := face.NewRecognizer(*flagDLIBModelsDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	recognizer = rec
}

func main() {
	defer recognizer.Close()

	e := echo.New()

	e.HideBanner = true

	e.POST("/:needle", func(c echo.Context) error {
		needle := strings.ToLower(c.Param("needle"))
		if needle != "all" && needle != "faces" {
			needle = "all"
		}

		file, err := c.FormFile("image")
		if err != nil {
			return c.JSON(400, map[string]interface{}{
				"success": false,
				"error":   "#1 - " + err.Error(),
			})
		}

		img, err := processSingleFileUpload(file)
		if err != nil {
			return c.JSON(400, map[string]interface{}{
				"success": false,
				"error":   "#2 - " + err.Error(),
			})
		}

		allProps := []image.Image{}

		if "faces" == needle {
			faces, err := processFacialExtraction(img)
			if err != nil {
				return c.JSON(500, map[string]interface{}{
					"success": false,
					"error":   "#3 - " + err.Error(),
				})
			}
			for _, face := range faces {
				allProps = append(allProps, getAllImageOrientations(face)...)
			}
		} else {
			allProps = append(allProps, getAllImageOrientations(img)...)
		}

		return c.JSON(200, map[string]interface{}{
			"success": true,
			"payload": processImagesHashes(allProps),
		})
	})

	e.Start(":4068")
}
