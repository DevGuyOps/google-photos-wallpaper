package main

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const (
	imagePath = "/home/guy/gphoto.jpg"
)

func main() {
	log.SetLevel(log.InfoLevel)

	if len(os.Args[:]) <= 1 {
		log.Error("Need more args")
	}

	// Auth user
	config := getClientConfig()
	client := getClient(config)

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
			getRandomPhotoFromAlbum(photoslibraryService, os.Args[3])
		case "list":
			getAlbumList(photoslibraryService)
		}
	}
}

func getPhotosFromAlbum(photoslibraryService *photoslibrary.Service, albumId string) {
	searchMediaRequest := &photoslibrary.SearchMediaItemsRequest{
		AlbumId:  albumId,
		PageSize: 50,
	}

	items, err := photoslibraryService.MediaItems.Search(searchMediaRequest).Do()
	if err != nil {
		log.Println(err)
	}

	numInAlbum := numItemsInAlbum(photoslibraryService, albumId)

	log.Println("here: " + string(random(0, int(numInAlbum))))

	for _, item := range items.MediaItems {
		downloadImage("items/"+item.Id+".jpg", item.BaseUrl, item.MediaMetadata.Width, item.MediaMetadata.Height)
	}
}

func getRandomPhotoFromAlbum(photoslibraryService *photoslibrary.Service, albumId string) {
	numInAlbum := numItemsInAlbum(photoslibraryService, albumId)
	photoNumInAlbum := random(0, numInAlbum)
	pageSize := 50
	pageCount := 0
	nextPageToken := ""

	log.Println("Rand Num: " + strconv.Itoa(photoNumInAlbum))

	for pageCount*pageSize < numInAlbum {
		var searchMediaRequest photoslibrary.SearchMediaItemsRequest

		// Get next page details
		if nextPageToken != "" {
			searchMediaRequest = photoslibrary.SearchMediaItemsRequest{
				AlbumId:   albumId,
				PageSize:  int64(pageSize),
				PageToken: nextPageToken,
				// Filters: &photoslibrary.Filters{
				// 	MediaTypeFilter: &photoslibrary.MediaTypeFilter {
				// 		MediaTypes: []string{
				// 			"PHOTO",
				// 		},
				// 	},
				// },
			}
		} else {
			searchMediaRequest = photoslibrary.SearchMediaItemsRequest{
				AlbumId:  albumId,
				PageSize: int64(pageSize),
				// Filters: &photoslibrary.Filters{
				// 	MediaTypeFilter: &photoslibrary.MediaTypeFilter {
				// 		MediaTypes: []string{
				// 			"PHOTO",
				// 		},
				// 	},
				// },
			}
		}

		// Perform search
		items, err := photoslibraryService.MediaItems.Search(&searchMediaRequest).Do()
		if err != nil {
			log.Println(err)
		}

		nextPageToken = items.NextPageToken

		// Download pic and do many things with it
		if pageCount*pageSize < photoNumInAlbum && pageCount*pageSize+pageSize > photoNumInAlbum {
			photoIndex := photoNumInAlbum - pageCount*pageSize
			currItem := items.MediaItems[photoIndex]

			downloadImage(imagePath, currItem.BaseUrl, currItem.MediaMetadata.Width, currItem.MediaMetadata.Height)
			setWallpaper(imagePath)
			break
		}

		pageCount++
	}
}

func getAlbumList(photoslibraryService *photoslibrary.Service) {
	albumList, err := photoslibraryService.Albums.List().Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, album := range albumList.Albums {
		log.Println(album.Title + " - [" + album.Id + "]")
	}
}

func numItemsInAlbum(photoslibraryService *photoslibrary.Service, albumId string) int {
	album, err := photoslibraryService.Albums.Get(albumId).Do()
	if err != nil {
		log.Fatal(err)
	}

	return int(album.MediaItemsCount)
}
