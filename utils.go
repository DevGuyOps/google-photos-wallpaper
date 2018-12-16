package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/reujab/wallpaper"
	log "github.com/sirupsen/logrus"
)

func downloadImage(filepath string, url string, width int64, height int64) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	fullURL := fmt.Sprintf("%s=w%s-h%s", url, strconv.FormatInt(width, 10), strconv.FormatInt(height, 10))
	resp, err := http.Get(fullURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Info("Photo downloaded")

	return nil
}

func setWallpaper(imagePath string) {
	wallpaper.SetFromFile(imagePath)
	log.Info("Wallpaper set")
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
