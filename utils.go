package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

	return nil
}
