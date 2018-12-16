package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	log "github.com/sirupsen/logrus"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

type Opts struct {
	ConfigPath flags.Filename `short:"c" long:"config" description:"Config file path" default:"config.json"`
}

func main() {
	if len(os.Args[:]) <= 1 {
		log.Error("Need more args")
	}

	var opts Opts
	flags.ParseArgs(&opts, os.Args[:])

	// Get config from file
	config, err := configFromFile(string(opts.ConfigPath))
	if err != nil {
		log.Error(err)
	}

	// Log setup
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)

	// Auth user
	clientConfig := getClientConfig(config.ClientConfPath)
	client := getClient(clientConfig)

	// Setup client
	photoslibraryService, err := photoslibrary.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// Deal with args
	switch strings.ToLower(os.Args[1]) {
	case "album":
		switch strings.ToLower(os.Args[2]) {
		case "random":
			mediaItem := getRandomPhotoFromAlbum(photoslibraryService, os.Args[3])
			if mediaItem != nil {
				downloadImage(config.WallpaperImgPath, mediaItem.BaseUrl, mediaItem.MediaMetadata.Width, mediaItem.MediaMetadata.Height)
				setWallpaper(config.WallpaperImgPath)
			}
		case "list":
			albumList := getAlbumList(photoslibraryService)
			for _, album := range albumList {
				log.Info(fmt.Sprintf("%s - [%s]", album.Title, album.Id))
			}
		}
	}
}
