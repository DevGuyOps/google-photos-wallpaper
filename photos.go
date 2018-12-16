package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

func downloadPhotosFromAlbum(photoslibraryService *photoslibrary.Service, albumId string) {
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

func getRandomPhotoFromAlbum(photoslibraryService *photoslibrary.Service, albumId string) *photoslibrary.MediaItem {
	numInAlbum := numItemsInAlbum(photoslibraryService, albumId)
	photoNumInAlbum := random(0, numInAlbum)
	pageSize := 50
	pageCount := 0
	nextPageToken := ""

	log.Debug("Rand Num: " + strconv.Itoa(photoNumInAlbum))

	for pageCount*pageSize < numInAlbum {
		var searchMediaRequest photoslibrary.SearchMediaItemsRequest

		// Get next page details
		if nextPageToken != "" {
			searchMediaRequest = photoslibrary.SearchMediaItemsRequest{
				AlbumId:   albumId,
				PageSize:  int64(pageSize),
				PageToken: nextPageToken,
				// Filters: &photoslibrary.Filters{
				// 	MediaTypeFilter: &photoslibrary.MediaTypeFilter{
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
				// 	MediaTypeFilter: &photoslibrary.MediaTypeFilter{
				// 		MediaTypes: []string{
				// 			"PHOTO",
				// 		},
				// 	},
				// },
			}
		}

		// Perform search
		log.Info("Searching for photo")
		items, err := photoslibraryService.MediaItems.Search(&searchMediaRequest).Do()
		if err != nil {
			log.Println(err)
		}

		nextPageToken = items.NextPageToken

		if pageCount*pageSize < photoNumInAlbum && pageCount*pageSize+pageSize > photoNumInAlbum {
			photoIndex := photoNumInAlbum - pageCount*pageSize
			currItem := items.MediaItems[photoIndex]

			return currItem
		}

		pageCount++
	}

	return nil
}

func getAlbumList(photoslibraryService *photoslibrary.Service) []*photoslibrary.Album {
	albumList, err := photoslibraryService.Albums.List().Do()
	if err != nil {
		log.Fatal(err)
	}

	return albumList.Albums
}

func numItemsInAlbum(photoslibraryService *photoslibrary.Service, albumId string) int {
	album, err := photoslibraryService.Albums.Get(albumId).Do()
	if err != nil {
		log.Fatal(err)
	}

	// Note: album.MediaItemsCount is a change made to the Google SDK in vendor
	// Google changed the API and have yet to change the SDK to bring it in line
	return int(album.MediaItemsCount)
}
