package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/rs/xid"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/labstack/echo"
)

var (
	flagHTTPAddr = flag.String("listen", ":4068", "the http listen address")
	flagDBPath   = flag.String("db", "./xmg-data", "the directory where the database will be located")
)

var (
	db *leveldb.DB
)

func init() {
	flag.Parse()
	var err error
	db, err = leveldb.OpenFile(*flagDBPath, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	defer db.Close()

	e := echo.New()
	e.HideBanner = true

	e.POST("/:action", func(c echo.Context) error {
		action := strings.ToLower(c.Param("action"))
		if action != "search" && action != "submit" {
			action = "search"
		}

		maxDistance, _ := strconv.Atoi(c.QueryParam("maxdist"))
		if maxDistance < 0 {
			maxDistance = 0
		}

		id := c.QueryParam("id")
		if "" == id {
			id = xid.New().String()
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

		allProps := getAllImageOrientations(img)
		hashes := processImagesHashes(allProps)

		if "submit" == action {
			if err := storeAppend(id, hashes...); err != nil {
				return c.JSON(400, map[string]interface{}{
					"success": false,
					"error":   "#3 - " + err.Error(),
				})
			}

			return c.JSON(200, map[string]interface{}{
				"success": true,
				"payload": id,
			})
		}

		return c.JSON(200, map[string]interface{}{
			"success": true,
			"payload": storeFind(maxDistance, hashes...),
		})
	})

	e.Start(":4068")
}
